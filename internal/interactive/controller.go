package interactive

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
	"golang.org/x/text/width"

	commandregistry "github.com/bmf-san/ggc/v7/cmd/command"
	"github.com/bmf-san/ggc/v7/internal/termio"
	"github.com/bmf-san/ggc/v7/pkg/config"
	"github.com/bmf-san/ggc/v7/pkg/git"
	uiutil "github.com/bmf-san/ggc/v7/pkg/ui"
)

// initialInputCapacity defines the initial capacity for the input rune buffer
// used by the real-time editor. It helps minimize reallocations during typing
// while keeping memory usage modest.
const initialInputCapacity = 64

// GitStatus represents the current Git repository status
type GitStatus struct {
	Branch     string
	Modified   int
	Staged     int
	Ahead      int
	Behind     int
	HasChanges bool
}

// ANSIColors is an alias to the shared UI palette definition.
type ANSIColors = uiutil.ANSIColors

// NewANSIColors exposes the shared ANSI color palette helper.
func NewANSIColors() *ANSIColors {
	return uiutil.NewANSIColors()
}

// getGitStatus retrieves the current Git repository status
func getGitStatus(gitClient git.StatusInfoReader) *GitStatus {
	status := &GitStatus{}

	// Get current branch name
	if branch := getGitBranch(gitClient); branch != "" {
		status.Branch = branch
	} else {
		return nil // Not in a git repository
	}

	// Get working directory status
	modified, staged := getGitWorkingStatus(gitClient)
	status.Modified = modified
	status.Staged = staged
	status.HasChanges = modified > 0 || staged > 0

	// Get remote tracking status
	ahead, behind := getGitRemoteStatus(gitClient)
	status.Ahead = ahead
	status.Behind = behind

	return status
}

// getGitBranch gets the current branch name
func getGitBranch(gitClient git.StatusInfoReader) string {
	branch, err := gitClient.GetCurrentBranch()
	if err != nil {
		return ""
	}
	return branch
}

// getGitWorkingStatus gets the number of modified and staged files
func getGitWorkingStatus(gitClient git.StatusInfoReader) (modified, staged int) {
	output, err := gitClient.StatusShortWithColor()
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		// First character: staged status
		// Second character: working tree status
		if line[0] != ' ' && line[0] != '?' {
			staged++
		}
		if line[1] != ' ' && line[1] != '?' {
			modified++
		}
	}
	return modified, staged
}

// getGitRemoteStatus gets ahead/behind count compared to remote
func getGitRemoteStatus(gitClient git.StatusInfoReader) (ahead, behind int) {
	output, err := gitClient.GetAheadBehindCount("HEAD", "@{upstream}")
	if err != nil {
		return 0, 0 // No upstream or other error
	}

	parts := strings.Fields(strings.TrimSpace(output))
	if len(parts) != 2 {
		return 0, 0
	}

	ahead, _ = strconv.Atoi(parts[0])
	behind, _ = strconv.Atoi(parts[1])
	return ahead, behind
}

// CommandInfo contains the name and description of the command
type CommandInfo struct {
	Command     string
	Description string
}

func buildInteractiveCommands() []CommandInfo {
	var list []CommandInfo
	allCommands := commandregistry.All()
	for i := range allCommands {
		cmd := &allCommands[i]
		if cmd.Hidden {
			continue
		}
		if len(cmd.Subcommands) == 0 {
			list = append(list, CommandInfo{Command: cmd.Name, Description: cmd.Summary})
			continue
		}
		for _, sub := range cmd.Subcommands {
			if sub.Hidden {
				continue
			}
			list = append(list, CommandInfo{Command: sub.Name, Description: sub.Summary})
		}
	}
	return list
}

// UI represents the interface for terminal UI operations
type UI struct {
	stdin                   io.Reader
	stdout                  io.Writer
	stderr                  io.Writer
	term                    termio.Terminal
	renderer                *Renderer
	state                   *UIState
	handler                 *KeyHandler
	colors                  *ANSIColors
	gitStatus               *GitStatus
	gitClient               git.StatusInfoReader
	reader                  *bufio.Reader
	contextMgr              *ContextManager
	profile                 Profile
	workflowMgr             *WorkflowManager
	workflow                *Workflow // Deprecated: maintained for legacy tests until refactor completes.
	workflowEx              *WorkflowExecutor
	softCancelFlash         atomic.Bool
	pendingWorkflowTemplate string
	pendingWorkflowCommand  string
	pendingWorkflowArgs     []string
	configMgr               *config.Manager
	config                  *config.Config
}

// UIState holds the current state of the interactive UI
type UIState struct {
	selected                int
	input                   string
	cursorPos               int // Cursor position in input string
	filtered                []CommandInfo
	context                 Context   // Current UI context (input/results/search/global)
	contextStack            []Context // Context stack for nested states
	onContextChange         func(Context, Context)
	showWorkflow            bool // Whether to show the workflow builder
	workflowSelectionActive bool
	workflowSelectionIndex  int
	workflowListIndex       int
}

// UpdateFiltered updates the filtered commands based on current input using fuzzy matching
func (s *UIState) UpdateFiltered() {
	input := strings.ToLower(s.input)
	if input == "" {
		s.filtered = make([]CommandInfo, len(commands))
		copy(s.filtered, commands)
	} else {
		type match struct {
			info  CommandInfo
			score matchScore
		}
		matches := make([]match, 0, len(commands))
		for _, cmd := range commands {
			cmdLower := strings.ToLower(cmd.Command)
			if ok, score := fuzzyMatchScore(cmdLower, input); ok {
				matches = append(matches, match{info: cmd, score: score})
			}
		}
		sort.SliceStable(matches, func(i, j int) bool {
			return matches[i].score.less(matches[j].score)
		})
		s.filtered = make([]CommandInfo, len(matches))
		for i, match := range matches {
			s.filtered[i] = match.info
		}
	}
	// Reset selection if out of bounds
	if s.selected >= len(s.filtered) {
		s.selected = len(s.filtered) - 1
	}
	if s.selected < 0 {
		s.selected = 0
	}
}

// Context Management Methods

// EnterContext pushes the current context onto the stack and switches to the new context
func (s *UIState) EnterContext(newContext Context) {
	if s.context == newContext {
		return
	}
	old := s.context
	s.contextStack = append(s.contextStack, s.context)
	s.context = newContext
	s.notifyContextChange(old, newContext)
}

// ExitContext pops the previous context from the stack
func (s *UIState) ExitContext() {
	if len(s.contextStack) > 0 {
		old := s.context
		s.context = s.contextStack[len(s.contextStack)-1]
		s.contextStack = s.contextStack[:len(s.contextStack)-1]
		s.notifyContextChange(old, s.context)
	} else if s.context != ContextGlobal {
		old := s.context
		s.context = ContextGlobal
		s.notifyContextChange(old, s.context)
	}
}

// GetCurrentContext returns the current UI context
func (s *UIState) GetCurrentContext() Context {
	return s.context
}

// SetContext directly sets the context (use with caution)
func (s *UIState) SetContext(ctx Context) {
	if s.context == ctx {
		return
	}
	old := s.context
	s.context = ctx
	s.notifyContextChange(old, ctx)
}

// notifyContextChange triggers the callback when the active context changes
func (s *UIState) notifyContextChange(oldCtx, newCtx Context) {
	if s.onContextChange != nil && oldCtx != newCtx {
		s.onContextChange(oldCtx, newCtx)
	}
}

// ActivateWorkflowSelection enables the workflow selection overlay.
func (s *UIState) ActivateWorkflowSelection(initialIndex int) {
	s.workflowSelectionActive = true
	s.workflowSelectionIndex = initialIndex
}

// DeactivateWorkflowSelection closes the workflow selection overlay.
func (s *UIState) DeactivateWorkflowSelection() {
	s.workflowSelectionActive = false
	s.workflowSelectionIndex = 0
}

// AdjustWorkflowSelection moves the selection pointer within bounds.
func (s *UIState) AdjustWorkflowSelection(delta, max int) {
	if !s.workflowSelectionActive || max <= 0 {
		return
	}
	s.workflowSelectionIndex += delta
	if s.workflowSelectionIndex < 0 {
		s.workflowSelectionIndex = max - 1
	} else if s.workflowSelectionIndex >= max {
		s.workflowSelectionIndex = 0
	}
}

// SetWorkflowListIndex updates the pointer used in the management view.
func (s *UIState) SetWorkflowListIndex(value, max int) {
	if max <= 0 {
		s.workflowListIndex = 0
		return
	}
	if value < 0 {
		value = 0
	} else if value >= max {
		value = max - 1
	}
	s.workflowListIndex = value
}

// IsInInputMode returns true if currently in input context
func (s *UIState) IsInInputMode() bool {
	return s.context == ContextInput
}

// IsInResultsMode returns true if currently in results context
func (s *UIState) IsInResultsMode() bool {
	return s.context == ContextResults
}

// IsInSearchMode returns true if currently in search context
func (s *UIState) IsInSearchMode() bool {
	return s.context == ContextSearch
}

// fuzzyMatch performs fuzzy matching between text and pattern
// Returns true if all characters in pattern appear in text in order (but not necessarily consecutive)
func fuzzyMatch(text, pattern string) bool {
	matched, _ := fuzzyMatchScore(text, pattern)
	return matched
}

// fuzzyMatchScore returns whether the pattern matches the text and a relevance score for sorting results.
// Lower scores indicate a tighter, earlier match.
func fuzzyMatchScore(text, pattern string) (bool, matchScore) {
	if pattern == "" {
		return true, matchScore{length: len([]rune(text))}
	}

	textRunes := []rune(text)
	patternRunes := []rune(pattern)

	matched, meta := matchPattern(textRunes, patternRunes)
	if !matched {
		return false, matchScore{}
	}

	trailing := len(textRunes) - meta.lastIndex - 1
	continuation := continuationPenalty(textRunes, meta.lastIndex)
	score := matchScore{
		first:        meta.firstIndex,
		gap:          meta.gapScore,
		trailing:     trailing,
		continuation: continuation,
		length:       len(textRunes),
	}

	return true, score
}

type matchMetadata struct {
	firstIndex int
	lastIndex  int
	gapScore   int
}

func matchPattern(textRunes, patternRunes []rune) (bool, matchMetadata) {
	meta := matchMetadata{
		firstIndex: -1,
		lastIndex:  -1,
	}

	textIdx := 0
	patternIdx := 0

	for textIdx < len(textRunes) && patternIdx < len(patternRunes) {
		if textRunes[textIdx] == patternRunes[patternIdx] {
			if meta.firstIndex == -1 {
				meta.firstIndex = textIdx
			}
			if meta.lastIndex != -1 {
				meta.gapScore += textIdx - meta.lastIndex - 1
			}
			meta.lastIndex = textIdx
			patternIdx++
		}
		textIdx++
	}

	if patternIdx != len(patternRunes) {
		return false, meta
	}

	return true, meta
}

func continuationPenalty(textRunes []rune, lastMatchIdx int) int {
	if lastMatchIdx < 0 || lastMatchIdx+1 >= len(textRunes) {
		return 0
	}

	nextIdx := lastMatchIdx + 1
	spaceSkipped := false
	for nextIdx < len(textRunes) && textRunes[nextIdx] == ' ' {
		spaceSkipped = true
		nextIdx++
	}

	if spaceSkipped && nextIdx < len(textRunes) && (unicode.IsLetter(textRunes[nextIdx]) || unicode.IsDigit(textRunes[nextIdx])) {
		return 1
	}

	return 0
}

type matchScore struct {
	first        int
	gap          int
	trailing     int
	continuation int
	length       int
}

func (m matchScore) less(other matchScore) bool {
	if m.first != other.first {
		return m.first < other.first
	}
	if m.gap != other.gap {
		return m.gap < other.gap
	}
	if m.continuation != other.continuation {
		return m.continuation < other.continuation
	}
	if m.trailing != other.trailing {
		return m.trailing < other.trailing
	}
	if m.length != other.length {
		return m.length < other.length
	}
	return false
}

// MoveUp moves selection up
func (s *UIState) MoveUp() {
	// Switch to results context when navigating
	if s.context != ContextResults && s.context != ContextSearch {
		s.SetContext(ContextResults)
	}

	if s.selected > 0 {
		s.selected--
	}
}

// MoveDown moves selection down
func (s *UIState) MoveDown() {
	// Switch to results context when navigating
	if s.context != ContextResults && s.context != ContextSearch {
		s.SetContext(ContextResults)
	}

	if s.selected < len(s.filtered)-1 {
		s.selected++
	}
}

// AddRune adds a UTF-8 rune to the input at cursor position
func (s *UIState) AddRune(r rune) {
	// Switch to input context when user starts typing
	if s.context != ContextInput {
		s.SetContext(ContextInput)
	}

	// Convert current input to runes for proper cursor positioning
	inputRunes := []rune(s.input)
	if s.cursorPos <= len(inputRunes) {
		// Efficiently insert the rune at the cursor position
		newRunes := make([]rune, len(inputRunes)+1)
		copy(newRunes, inputRunes[:s.cursorPos])
		newRunes[s.cursorPos] = r
		copy(newRunes[s.cursorPos+1:], inputRunes[s.cursorPos:])

		s.input = string(newRunes)
		s.cursorPos++
		s.UpdateFiltered()

		// Switch to search context when actively filtering
		if s.input != "" && s.context != ContextSearch {
			s.SetContext(ContextSearch)
		}
	}
}

// RemoveChar removes character before cursor (backspace)
func (s *UIState) RemoveChar() {
	if s.cursorPos > 0 && s.input != "" {
		// Convert to runes for proper UTF-8 handling
		inputRunes := []rune(s.input)
		if s.cursorPos <= len(inputRunes) {
			inputRunes = append(inputRunes[:s.cursorPos-1], inputRunes[s.cursorPos:]...)

			s.input = string(inputRunes)
			s.cursorPos--
			s.UpdateFiltered()
		}
	}
}

// ClearInput clears all input
func (s *UIState) ClearInput() {
	s.input = ""
	s.cursorPos = 0
	s.UpdateFiltered()
}

// DeleteWord deletes word before cursor (Ctrl+W)
func (s *UIState) DeleteWord() {
	if s.cursorPos == 0 {
		return
	}

	// Convert to runes for proper UTF-8 handling
	inputRunes := []rune(s.input)

	// Find start of current word (skip trailing spaces first)
	pos := s.cursorPos - 1
	for pos >= 0 && inputRunes[pos] == ' ' {
		pos--
	}

	// Find start of word
	for pos >= 0 && inputRunes[pos] != ' ' {
		pos--
	}
	pos++ // Move to first character of word

	// Delete from word start to cursor
	inputRunes = append(inputRunes[:pos], inputRunes[s.cursorPos:]...)
	s.input = string(inputRunes)
	s.cursorPos = pos
	s.UpdateFiltered()
}

// DeleteToEnd deletes from cursor to end of line (Ctrl+K)
func (s *UIState) DeleteToEnd() {
	if s.cursorPos < utf8.RuneCountInString(s.input) {
		inputRunes := []rune(s.input)
		s.input = string(inputRunes[:s.cursorPos])
		s.UpdateFiltered()
	}
}

// MoveToBeginning moves cursor to beginning of line (Ctrl+A)
func (s *UIState) MoveToBeginning() {
	s.cursorPos = 0
}

// MoveToEnd moves cursor to end of line (Ctrl+E)
func (s *UIState) MoveToEnd() {
	s.cursorPos = utf8.RuneCountInString(s.input)
}

// MoveLeft moves cursor one rune left
func (s *UIState) MoveLeft() {
	if s.cursorPos > 0 {
		s.cursorPos--
	}
}

// MoveRight moves cursor one rune right
func (s *UIState) MoveRight() {
	if s.cursorPos < utf8.RuneCountInString(s.input) {
		s.cursorPos++
	}
}

// MoveWordLeft moves cursor to the beginning of the previous word
func (s *UIState) MoveWordLeft() {
	if s.cursorPos == 0 {
		return
	}
	runes := []rune(s.input)
	pos := s.cursorPos - 1
	// Skip any spaces to the left
	for pos >= 0 && unicode.IsSpace(runes[pos]) {
		pos--
	}
	// Skip the word characters to the left
	for pos >= 0 && !unicode.IsSpace(runes[pos]) {
		pos--
	}
	s.cursorPos = pos + 1
}

// MoveWordRight moves cursor to the start of the next word
func (s *UIState) MoveWordRight() {
	runes := []rune(s.input)
	n := len(runes)
	pos := s.cursorPos
	if pos >= n {
		return
	}
	// Skip current word characters
	for pos < n && !unicode.IsSpace(runes[pos]) {
		pos++
	}
	// Skip following spaces
	for pos < n && unicode.IsSpace(runes[pos]) {
		pos++
	}
	s.cursorPos = pos
}

// GetSelectedCommand returns the currently selected command
func (s *UIState) GetSelectedCommand() *CommandInfo {
	if len(s.filtered) > 0 && s.selected >= 0 && s.selected < len(s.filtered) {
		return &s.filtered[s.selected]
	}
	return nil
}

// HasInput returns true if there is input
func (s *UIState) HasInput() bool {
	return s.input != ""
}

// HasMatches returns true if there are filtered matches
func (s *UIState) HasMatches() bool {
	return len(s.filtered) > 0
}

// Renderer handles all terminal rendering operations
type Renderer struct {
	writer io.Writer
	width  int
	height int
	colors *ANSIColors
}

type keybindHelpEntry struct {
	key  string
	desc string
}

// KeyHandler manages keyboard input processing
type KeyHandler struct {
	ui            *UI
	contextualMap *ContextualKeyBindingMap
}

// GetCurrentKeyMap returns the appropriate keybinding map for the current context
func (h *KeyHandler) GetCurrentKeyMap() *KeyBindingMap {
	if h == nil {
		return DefaultKeyBindingMap()
	}
	if h.contextualMap != nil && h.ui != nil && h.ui.state != nil {
		currentContext := h.ui.state.GetCurrentContext()
		if contextMap, exists := h.contextualMap.GetContext(currentContext); exists && contextMap != nil {
			return contextMap
		}
		if contextMap, exists := h.contextualMap.GetContext(ContextGlobal); exists && contextMap != nil {
			return contextMap
		}
	}

	return DefaultKeyBindingMap()
}

// HandleKey processes UTF-8 rune input and returns true if should continue
// This method handles both single-byte (ASCII/control) and multibyte characters
func (h *KeyHandler) HandleKey(r rune, _ bool, oldState *term.State, reader *bufio.Reader) (bool, []string) {
	// Set the reader for consistent access during escape sequence handling
	h.ui.reader = reader
	// Handle workflow-specific keys first (Tab, etc.)
	if handled := h.handleWorkflowKeys(r); handled {
		return true, nil
	}

	// Handle control characters (ASCII range)
	if r < 128 && unicode.IsControl(r) { // ASCII control characters
		if handled, shouldContinue, result := h.handleControlChar(byte(r), oldState, reader); handled {
			return shouldContinue, result
		}
	}

	// Handle printable characters (both ASCII and multibyte)
	if unicode.IsPrint(r) {
		// Don't accept text input in workflow view
		if !h.ui.state.showWorkflow {
			h.ui.state.AddRune(r)
		}
	}
	return true, nil
}

// handleWorkflowKeys processes workflow-related key bindings and returns (handled, result)
func (h *KeyHandler) handleWorkflowKeys(r rune) bool {
	km := h.GetCurrentKeyMap()
	if km == nil {
		km = DefaultKeyBindingMap()
	}

	keyStroke := NewCharKeyStroke(r)
	state := h.ui.state

	if state.workflowSelectionActive {
		return h.handleWorkflowSelectionKey(km, keyStroke)
	}

	if km.MatchesKeyStroke("workflow_cancel", keyStroke) {
		if h.handleWorkflowCancel(state) {
			return true
		}
	}

	if !state.showWorkflow && km.MatchesKeyStroke("clear_workflow", keyStroke) {
		h.clearWorkflow()
		return true
	}

	if !state.showWorkflow {
		return h.handleSearchWorkflowKey(km, keyStroke, state)
	}

	return h.handleWorkflowKeysInView(km, keyStroke)
}

func (h *KeyHandler) handleSearchWorkflowKey(km *KeyBindingMap, keyStroke KeyStroke, state *UIState) bool {
	if !km.MatchesKeyStroke("add_to_workflow", keyStroke) {
		return false
	}
	if state.HasMatches() {
		if selectedCmd := state.GetSelectedCommand(); selectedCmd != nil {
			h.addCommandToWorkflow(selectedCmd.Command)
		}
	}
	return true
}

func (h *KeyHandler) handleWorkflowSelectionKey(km *KeyBindingMap, keyStroke KeyStroke) bool {
	switch {
	case km.MatchesKeyStroke("add_to_workflow", keyStroke):
		h.moveWorkflowSelection(1)
	case km.MatchesKeyStroke("workflow_cancel", keyStroke):
		h.cancelWorkflowSelection()
	}
	return true
}

func (h *KeyHandler) handleWorkflowKeysInView(km *KeyBindingMap, keyStroke KeyStroke) bool {
	switch {
	case km.MatchesKeyStroke("workflow_create", keyStroke):
		h.createWorkflowFromView()
		return true
	case km.MatchesKeyStroke("workflow_copy", keyStroke):
		h.copyWorkflowFromView()
		return true
	case km.MatchesKeyStroke("workflow_save", keyStroke):
		h.saveWorkflowFromView()
		return true
	case km.MatchesKeyStroke("workflow_delete", keyStroke) || km.MatchesKeyStroke("clear_workflow", keyStroke):
		h.clearWorkflow()
		return true
	case km.MatchesKeyStroke("add_to_workflow", keyStroke):
		return true
	default:
		return false
	}
}

func (h *KeyHandler) handleWorkflowCancel(state *UIState) bool {
	switch {
	case state.workflowSelectionActive:
		h.cancelWorkflowSelection()
		return true
	case state.showWorkflow:
		h.ui.ToggleWorkflowView()
		return true
	default:
		return false
	}
}

// handleControlChar processes control characters and returns (handled, shouldContinue, result)
// Enhanced to support KeyStroke matching while maintaining backward compatibility
//
//nolint:revive // Control character handling inherently requires many cases
func (h *KeyHandler) handleControlChar(b byte, oldState *term.State, reader *bufio.Reader) (bool, bool, []string) {
	// Get the appropriate keybinding map for current context
	km := h.GetCurrentKeyMap()
	state := h.ui.state

	if state.workflowSelectionActive {
		if b >= 1 && b <= 26 {
			ctrlStroke := NewCtrlKeyStroke(rune('a' + b - 1))
			if km.MatchesKeyStroke("move_up", ctrlStroke) {
				h.moveWorkflowSelection(-1)
				return true, true, nil
			}
			if km.MatchesKeyStroke("move_down", ctrlStroke) {
				h.moveWorkflowSelection(1)
				return true, true, nil
			}
			if km.MatchesKeyStroke("soft_cancel", ctrlStroke) {
				h.cancelWorkflowSelection()
				return true, true, nil
			}
			if km.MatchesKeyStroke("workflow_cancel", ctrlStroke) {
				h.cancelWorkflowSelection()
				return true, true, nil
			}
		}

		switch b {
		case 13: // Enter
			h.finalizeWorkflowSelection()
			return true, true, nil
		case 27: // Esc
			h.cancelWorkflowSelection()
			return true, true, nil
		default:
			return true, true, nil
		}
	}

	// Create KeyStroke for this control character
	if b >= 1 && b <= 26 {
		// Control character: convert back to letter
		ctrlStroke := NewCtrlKeyStroke(rune('a' + b - 1))

		// Check each action using new KeyStroke matching
		if km.MatchesKeyStroke("move_up", ctrlStroke) {
			switch {
			case state.showWorkflow:
				h.moveWorkflowList(-1)
			default:
				h.ui.state.MoveUp()
			}
			return true, true, nil
		}
		if km.MatchesKeyStroke("move_down", ctrlStroke) {
			switch {
			case state.showWorkflow:
				h.moveWorkflowList(1)
			default:
				h.ui.state.MoveDown()
			}
			return true, true, nil
		}
		if km.MatchesKeyStroke("clear_line", ctrlStroke) {
			h.ui.state.ClearInput()
			return true, true, nil
		}
		if km.MatchesKeyStroke("delete_word", ctrlStroke) {
			h.ui.state.DeleteWord()
			return true, true, nil
		}
		if km.MatchesKeyStroke("delete_to_end", ctrlStroke) {
			h.ui.state.DeleteToEnd()
			return true, true, nil
		}
		if km.MatchesKeyStroke("move_to_beginning", ctrlStroke) {
			h.ui.state.MoveToBeginning()
			return true, true, nil
		}
		if km.MatchesKeyStroke("move_to_end", ctrlStroke) {
			h.ui.state.MoveToEnd()
			return true, true, nil
		}

		// Check for workflow toggle
		if km.MatchesKeyStroke("toggle_workflow_view", ctrlStroke) && h.ui.state.input == "" {
			h.ui.ToggleWorkflowView()
			return true, true, nil
		}
		if km.MatchesKeyStroke("workflow_create", ctrlStroke) {
			if state.showWorkflow {
				h.createWorkflowFromView()
			}
			return true, true, nil
		}
		if km.MatchesKeyStroke("workflow_cancel", ctrlStroke) {
			if state.showWorkflow {
				h.ui.ToggleWorkflowView()
			}
			return true, true, nil
		}
		if km.MatchesKeyStroke("soft_cancel", ctrlStroke) {
			h.handleSoftCancel(oldState)
			return true, true, nil
		}
	}

	// Handle special cases that are not Ctrl+letter
	switch b {
	case 3: // Ctrl+C
		h.handleCtrlC(oldState)
		return true, false, nil
	case 13: // Enter
		if state.workflowSelectionActive {
			h.finalizeWorkflowSelection()
			return true, true, nil
		}
		shouldContinue, result := h.handleEnter(oldState)
		return true, shouldContinue, result
	case 127, 8: // Backspace
		h.ui.state.RemoveChar()
		return true, true, nil
	case 27: // ESC: arrow keys and Option/Alt modifiers
		escStroke := NewEscapeKeyStroke()
		if km.MatchesKeyStroke("workflow_cancel", escStroke) {
			if h.handleWorkflowCancel(state) {
				return true, true, nil
			}
		}
		if h.shouldHandleEscapeAsSoftCancel() {
			h.handleSoftCancel(oldState)
			return true, true, nil
		}
		h.handleEscapeSequence(reader)
		return true, true, nil
	default:
		return false, true, nil
	}
}

// handleEscapeSequence parses common ESC sequences for arrow and word navigation.
// Supports:
// - Arrow keys: ESC [ C/D (right/left), ESC O C/D (application mode)
// - Ctrl+Arrow: ESC [ 1;5 C/D or ESC [ 5 C/D
// - Alt/Option+Arrow: ESC [ 1;3 C/D, ESC [ 1;9 C/D (varies by terminal)
// - macOS Option word nav: ESC b / ESC f
func (h *KeyHandler) handleEscapeSequence(reader *bufio.Reader) {
	if h.ui == nil {
		return
	}

	// Read next byte after ESC
	var b byte
	var err error

	if reader != nil {
		// Use provided buffered reader (non-raw mode)
		b, err = reader.ReadByte()
	} else {
		// Raw mode: read directly from stdin
		var buf [1]byte
		_, err = h.ui.stdin.Read(buf[:])
		b = buf[0]
	}

	if err != nil {
		return
	}

	switch b {
	case '[':
		h.handleCSISequence(reader)
	case 'O':
		h.handleApplicationCursorMode(reader)
	case 'b':
		h.ui.state.MoveWordLeft()
	case 'f':
		h.ui.state.MoveWordRight()
	case 127, 8:
		// Meta-Backspace (Option+Backspace): delete word left
		h.ui.state.DeleteWord()
	}
}

func (h *KeyHandler) handleSoftCancel(_ *term.State) {
	if h == nil || h.ui == nil {
		return
	}

	if h.ui.resetToSearchMode() {
		h.ui.notifySoftCancel()
	}
}

func (h *KeyHandler) shouldHandleEscapeAsSoftCancel() bool {
	km := h.GetCurrentKeyMap()
	if km == nil || !km.MatchesKeyStroke("soft_cancel", NewEscapeKeyStroke()) {
		return false
	}

	if h.ui == nil {
		return false
	}

	if h.ui.reader != nil && h.ui.reader.Buffered() > 0 {
		return false
	}

	if file, ok := h.ui.stdin.(*os.File); ok {
		if pending, err := termio.PendingInput(file.Fd()); err == nil {
			return pending == 0
		}
	}

	return false
}

// handleCSISequence handles CSI (Control Sequence Introducer) sequences
func (h *KeyHandler) handleCSISequence(reader *bufio.Reader) {
	var params []byte
	for {
		var nb byte
		var err error

		if reader != nil {
			// Use provided buffered reader (non-raw mode)
			nb, err = reader.ReadByte()
		} else {
			// Raw mode: read directly from stdin
			var buf [1]byte
			_, err = h.ui.stdin.Read(buf[:])
			nb = buf[0]
		}

		if err != nil {
			return
		}
		if (nb >= 'A' && nb <= 'Z') || nb == '~' {
			h.processCSIFinalByte(nb, string(params))
			return
		}
		params = append(params, nb)
	}
}

// isWordMotionParam reports whether CSI params include a word-motion modifier
// commonly emitted by terminals (e.g., Ctrl/Alt variants use 5/3/9).
func isWordMotionParam(params string) bool {
	return strings.Contains(params, "5") || strings.Contains(params, "3") || strings.Contains(params, "9")
}

// processCSIFinalByte processes the final byte of a CSI sequence
func (h *KeyHandler) processCSIFinalByte(final byte, params string) {
	isWord := isWordMotionParam(params)

	// Build the full escape sequence for keybinding matching
	seq := h.buildCSISequence(final, params)
	keyStroke := NewRawKeyStroke(seq)
	km := h.GetCurrentKeyMap()

	// Try keybinding-based handling first
	if h.tryArrowKeybinding(km, keyStroke) {
		return
	}

	// Fallback to default cursor movement and word navigation
	h.handleDefaultArrowMovement(final, isWord)
}

// buildCSISequence builds a CSI escape sequence
func (h *KeyHandler) buildCSISequence(final byte, params string) []byte {
	if params == "" {
		return []byte{27, '[', final}
	}
	seq := append([]byte{27, '['}, []byte(params)...)
	return append(seq, final)
}

// tryArrowKeybinding attempts to handle arrow keys via keybindings
func (h *KeyHandler) tryArrowKeybinding(km *KeyBindingMap, keyStroke KeyStroke) bool {
	if km.MatchesKeyStroke("move_up", keyStroke) {
		if !h.ui.state.showWorkflow {
			h.ui.state.MoveUp()
		}
		return true
	}
	if km.MatchesKeyStroke("move_down", keyStroke) {
		if !h.ui.state.showWorkflow {
			h.ui.state.MoveDown()
		}
		return true
	}
	if km.MatchesKeyStroke("move_left", keyStroke) {
		h.ui.state.MoveLeft()
		return true
	}
	if km.MatchesKeyStroke("move_right", keyStroke) {
		h.ui.state.MoveRight()
		return true
	}
	return false
}

// handleDefaultArrowMovement handles default arrow key behavior
func (h *KeyHandler) handleDefaultArrowMovement(final byte, isWord bool) {
	switch final {
	case 'C': // Right
		if isWord {
			h.ui.state.MoveWordRight()
		} else {
			h.ui.state.MoveRight()
		}
	case 'D': // Left
		if isWord {
			h.ui.state.MoveWordLeft()
		} else {
			h.ui.state.MoveLeft()
		}
	}
}

// handleApplicationCursorMode handles application cursor mode sequences
func (h *KeyHandler) handleApplicationCursorMode(reader *bufio.Reader) {
	nb, err := h.readNextByte(reader)
	if err != nil {
		return
	}

	// Build the full escape sequence: ESC O <final>
	seq := []byte{27, 'O', nb}
	keyStroke := NewRawKeyStroke(seq)
	km := h.GetCurrentKeyMap()

	// Try keybinding-based handling first
	if h.tryArrowKeybinding(km, keyStroke) {
		return
	}

	// Fallback to default arrow key behavior
	h.handleDefaultAppCursorMovement(nb)
}

// handleDefaultAppCursorMovement handles default application cursor mode arrow keys
func (h *KeyHandler) handleDefaultAppCursorMovement(nb byte) {
	switch nb {
	case 'A':
		h.handleVerticalNavigation(-1)
	case 'B':
		h.handleVerticalNavigation(1)
	case 'C':
		h.ui.state.MoveRight()
	case 'D':
		h.ui.state.MoveLeft()
	}
}

func (h *KeyHandler) handleVerticalNavigation(delta int) {
	switch {
	case h.ui.state.workflowSelectionActive:
		h.moveWorkflowSelection(delta)
	case h.ui.state.showWorkflow:
		h.moveWorkflowList(delta)
	case delta < 0:
		h.ui.state.MoveUp()
	default:
		h.ui.state.MoveDown()
	}
}

// readNextByte reads the next byte from either a buffered reader or stdin
func (h *KeyHandler) readNextByte(reader *bufio.Reader) (byte, error) {
	if reader != nil {
		return reader.ReadByte()
	}
	var buf [1]byte
	_, err := h.ui.stdin.Read(buf[:])
	return buf[0], err
}

// handleCtrlC handles Ctrl+C key press
func (h *KeyHandler) handleCtrlC(oldState *term.State) {
	if oldState != nil {
		if f, ok := h.ui.stdin.(*os.File); ok {
			if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
				h.ui.writeError("failed to restore terminal state: %v", err)
			}
		}
	}
	h.ui.writeln("\nExiting...")
	os.Exit(0)
}

// restoreTerminalState restores the terminal from raw mode to cooked mode
func (h *KeyHandler) restoreTerminalState(oldState *term.State) {
	if oldState == nil {
		return
	}
	if f, ok := h.ui.stdin.(*os.File); ok {
		if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
			h.ui.writeError("failed to restore terminal state: %v", err)
		}
	}
}

// reenterRawMode re-enters raw mode after being restored
func (h *KeyHandler) reenterRawMode(oldState *term.State) {
	if oldState == nil {
		return
	}
	if f, ok := h.ui.stdin.(*os.File); ok {
		fd := int(f.Fd())
		if _, err := h.ui.term.MakeRaw(fd); err != nil {
			h.ui.writeError("failed to set terminal to raw mode: %v", err)
		}
	}
}

// handleEnter handles Enter key press
func (h *KeyHandler) handleEnter(oldState *term.State) (bool, []string) {
	// Handle workflow mode
	if h.ui.state.showWorkflow {
		// Execute workflow
		return h.executeWorkflow(oldState)
	}

	if !h.ui.state.HasInput() {
		return true, nil
	}

	selectedCmd := h.ui.state.GetSelectedCommand()
	if selectedCmd == nil {
		return true, nil
	}

	// Restore terminal state BEFORE showing Execute message
	h.restoreTerminalState(oldState)

	// Clear screen and show execution message
	clearScreen(h.ui.stdout)
	executeMsg := fmt.Sprintf("%s🚀 %sExecuting:%s %s%s%s\n\n",
		h.ui.colors.BrightGreen,
		h.ui.colors.BrightWhite+h.ui.colors.Bold,
		h.ui.colors.Reset,
		h.ui.colors.BrightCyan+h.ui.colors.Bold,
		selectedCmd.Command,
		h.ui.colors.Reset)
	h.ui.writeColor(executeMsg)

	// Handle placeholders
	args, canceled := h.processCommand(selectedCmd.Command)
	if canceled {
		// Re-enter raw mode before returning to main loop
		h.reenterRawMode(oldState)
		return true, nil
	}
	return false, args
}

// processCommand processes the command with placeholder replacement
func (h *KeyHandler) processCommand(cmdTemplate string) ([]string, bool) {
	placeholders := extractPlaceholders(cmdTemplate)

	if len(placeholders) == 0 {
		// No placeholders - execute immediately
		args := []string{"ggc"}
		args = append(args, strings.Fields(cmdTemplate)...)
		return args, false
	}

	// Interactive input for placeholders
	inputs, canceled := h.interactiveInput(placeholders)
	if canceled {
		h.handleSoftCancel(nil)
		return nil, true
	}

	// Placeholder replacement
	finalCmd := cmdTemplate
	for ph, val := range inputs {
		finalCmd = strings.ReplaceAll(finalCmd, "<"+ph+">", val)
	}

	args := []string{"ggc"}
	args = append(args, strings.Fields(finalCmd)...)
	return args, false
}

// interactiveInput provides real-time interactive input for placeholders
func (h *KeyHandler) interactiveInput(placeholders []string) (map[string]string, bool) {
	inputs := make(map[string]string)

	for i, ph := range placeholders {
		h.ui.write("\n")

		// Show progress and prompt
		if len(placeholders) > 1 {
			h.ui.write("%s[%d/%d]%s ",
				h.ui.colors.BrightBlue+h.ui.colors.Bold,
				i+1, len(placeholders),
				h.ui.colors.Reset)
		}

		h.ui.write("%s? %s%s%s: ",
			h.ui.colors.BrightGreen,
			h.ui.colors.BrightWhite+h.ui.colors.Bold,
			ph,
			h.ui.colors.Reset)

		// Get input with real-time feedback
		value, canceled := h.ui.readPlaceholderInput()
		if canceled {
			return nil, true
		}
		inputs[ph] = value

		// Show confirmation
		h.ui.write("%s✓ %s%s: %s%s%s\n",
			h.ui.colors.BrightGreen,
			h.ui.colors.BrightBlue,
			ph,
			h.ui.colors.BrightYellow+h.ui.colors.Bold,
			value,
			h.ui.colors.Reset)
	}

	return inputs, false
}

// getRealTimeInput gets user input with real-time display using raw terminal mode
func (h *KeyHandler) getRealTimeInput() (string, bool) {
	fd := int(os.Stdin.Fd())
	oldState, err := h.ui.term.MakeRaw(fd)
	if err != nil {
		return h.getLineInput()
	}
	defer func() { _ = h.ui.term.Restore(fd, oldState) }()

	return h.processRealTimeInput()
}

// processRealTimeInput handles the main input processing loop
func (h *KeyHandler) processRealTimeInput() (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	inputRunes := make([]rune, 0, initialInputCapacity)
	cursor := 0

	editor := &realTimeEditor{
		ui:         h.ui,
		inputRunes: &inputRunes,
		cursor:     &cursor,
	}

	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			continue
		}

		result := editor.handleInput(r, reader)
		if result.done {
			return result.text, false
		}
		if result.canceled {
			return "", true
		}
	}
	return string(inputRunes), false
}

// inputResult represents the result of handling input
type inputResult struct {
	done     bool
	canceled bool
	text     string
}

// realTimeEditor handles real-time input editing
type realTimeEditor struct {
	ui         *UI
	inputRunes *[]rune
	cursor     *int
}

// handleInput processes a single input rune
func (e *realTimeEditor) handleInput(r rune, reader *bufio.Reader) inputResult {
	switch r {
	case '\n', '\r':
		return e.handleEnter()
	case 3: // Ctrl+C
		return e.handleCtrlC()
	case 7: // Ctrl+G
		return e.handleSoftCancel()
	case 127, '\b': // Backspace
		e.handleBackspace()
		return inputResult{}
	case 27: // ESC sequences
		if e.shouldSoftCancelOnEscape(reader) {
			return e.handleSoftCancel()
		}
		e.handleEscape(reader)
		return inputResult{}
	default:
		if unicode.IsPrint(r) {
			e.handlePrintableChar(r)
		}
		return inputResult{}
	}
}

// handleEnter processes Enter key
func (e *realTimeEditor) handleEnter() inputResult {
	if len(*e.inputRunes) > 0 {
		e.ui.write("\r\n")
		return inputResult{done: true, text: string(*e.inputRunes)}
	}
	e.ui.write(" %s(required)%s", e.ui.colors.BrightRed, e.ui.colors.Reset)
	return inputResult{}
}

// handleCtrlC processes Ctrl+C
func (e *realTimeEditor) handleCtrlC() inputResult {
	e.ui.write("\r\n%sOperation canceled%s\r\n", e.ui.colors.BrightRed, e.ui.colors.Reset)
	return inputResult{canceled: true}
}

func (e *realTimeEditor) handleSoftCancel() inputResult {
	e.ui.write("\r\n")
	return inputResult{canceled: true}
}

// handleBackspace processes backspace key
func (e *realTimeEditor) handleBackspace() {
	if *e.cursor == 0 {
		return
	}
	start := e.findGraphemeStart(*e.cursor - 1)
	// Compute columns to move left/clear for the removed cluster
	moveCols := e.colsBetween(start, *e.cursor)
	clearedCols := 0
	for i := start; i < *e.cursor; i++ {
		clearedCols += e.runeWidth((*e.inputRunes)[i])
	}
	// Move cursor left, remove runes, and redraw tail
	e.moveLeft(moveCols)
	*e.inputRunes = append((*e.inputRunes)[:start], (*e.inputRunes)[*e.cursor:]...)
	*e.cursor = start
	e.printTailAndReposition(*e.cursor, clearedCols)
}

// findGraphemeStart finds the start of the grapheme cluster ending at the given position
func (e *realTimeEditor) findGraphemeStart(pos int) int {
	start := pos
	start = e.skipCombiningMarks(start)
	start = e.handleRegionalIndicators(start)
	start = e.handleZWJSequences(start)
	if start < 0 {
		start = 0
	}
	return start
}

// skipCombiningMarks skips any trailing variation selectors or combining marks
func (e *realTimeEditor) skipCombiningMarks(start int) int {
	for start >= 0 && (isCombining((*e.inputRunes)[start]) || isVariationSelector((*e.inputRunes)[start])) {
		start--
	}
	return start
}

// handleRegionalIndicators handles regional indicator pairs (flags)
func (e *realTimeEditor) handleRegionalIndicators(start int) int {
	if start >= 0 && isRegionalIndicator((*e.inputRunes)[start]) {
		if start > 0 && isRegionalIndicator((*e.inputRunes)[start-1]) {
			start--
		}
	}
	return start
}

// handleZWJSequences handles ZWJ sequences by including the joiner and previous rune repeatedly
func (e *realTimeEditor) handleZWJSequences(start int) int {
	for {
		if start > 0 && isZWJ((*e.inputRunes)[start-1]) {
			// Include ZWJ and the previous rune
			start -= 2
			// Also include any combining marks attached to the new base
			start = e.skipCombiningMarks(start)
			continue
		}
		break
	}
	return start
}

// handlePrintableChar processes printable characters
func (e *realTimeEditor) handlePrintableChar(r rune) {
	if *e.cursor == len(*e.inputRunes) {
		*e.inputRunes = append(*e.inputRunes, r)
	} else {
		*e.inputRunes = append((*e.inputRunes)[:*e.cursor], append([]rune{r}, (*e.inputRunes)[*e.cursor:]...)...)
	}
	e.ui.write("%s", string(r))
	*e.cursor++
	e.printTailAndReposition(*e.cursor, 0)
}

// handleEscape processes escape sequences for real-time input
func (e *realTimeEditor) handleEscape(reader *bufio.Reader) {
	b, err := reader.ReadByte()
	if err != nil {
		return
	}
	switch b {
	case '[':
		e.handleCSIEscape(reader)
	case 'O':
		e.handleApplicationEscape(reader)
	case 'b':
		e.moveWordLeft()
	case 'f':
		e.moveWordRight()
	case 127, '\b':
		// Option+Backspace: delete previous word
		e.deleteWordLeft()
	}
}

func (e *realTimeEditor) shouldSoftCancelOnEscape(reader *bufio.Reader) bool {
	if reader != nil && reader.Buffered() > 0 {
		return false
	}

	if os.Stdin != nil {
		if pending, err := termio.PendingInput(os.Stdin.Fd()); err == nil {
			return pending == 0
		}
	}

	return false
}

// handleCSIEscape processes CSI escape sequences for real-time input
func (e *realTimeEditor) handleCSIEscape(reader *bufio.Reader) {
	var params []byte
	for {
		nb, err := reader.ReadByte()
		if err != nil {
			return
		}
		if (nb >= 'A' && nb <= 'Z') || nb == '~' {
			e.processCSIEscape(nb, string(params))
			return
		}
		params = append(params, nb)
	}
}

// processCSIEscape handles CSI final byte for real-time input
func (e *realTimeEditor) processCSIEscape(final byte, params string) {
	isWord := isWordMotionParam(params)
	switch final {
	case 'C': // Right
		if isWord {
			e.moveWordRight()
		} else if *e.cursor < len(*e.inputRunes) {
			e.moveRight(e.runeWidth((*e.inputRunes)[*e.cursor]))
			*e.cursor++
		}
	case 'D': // Left
		if isWord {
			e.moveWordLeft()
		} else if *e.cursor > 0 {
			e.moveLeft(e.runeWidth((*e.inputRunes)[*e.cursor-1]))
			*e.cursor--
		}
	}
}

// handleApplicationEscape processes application mode escape sequences
func (e *realTimeEditor) handleApplicationEscape(reader *bufio.Reader) {
	nb, err := reader.ReadByte()
	if err != nil {
		return
	}
	switch nb {
	case 'C':
		if *e.cursor < len(*e.inputRunes) {
			e.moveRight(e.runeWidth((*e.inputRunes)[*e.cursor]))
			*e.cursor++
		}
	case 'D':
		if *e.cursor > 0 {
			e.moveLeft(e.runeWidth((*e.inputRunes)[*e.cursor-1]))
			*e.cursor--
		}
	}
}

// Helper methods for realTimeEditor

// isCombining reports whether r is a combining mark (zero display width)
func isCombining(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}

// isVariationSelector reports whether r is a variation selector (zero width)
func isVariationSelector(r rune) bool {
	// U+FE00..U+FE0F (VS1..VS16) and U+E0100..U+E01EF (IVS)
	return (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF)
}

// isRegionalIndicator reports whether r is a regional indicator rune (used for flags)
func isRegionalIndicator(r rune) bool { return r >= 0x1F1E6 && r <= 0x1F1FF }

// isZWJ reports whether r is ZERO WIDTH JOINER
func isZWJ(r rune) bool { return r == 0x200D }

// isEmoji reports common emoji ranges that should render as width 2 on most terminals
func isEmoji(r rune) bool {
	return isEmojiRange1(r) || isEmojiRange2(r)
}

// isEmojiRange1 checks the first set of emoji Unicode ranges
func isEmojiRange1(r rune) bool {
	return (r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols
		(r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		(r >= 0x1F780 && r <= 0x1F7FF) // Geometric Shapes Extended
}

// isEmojiRange2 checks the second set of emoji Unicode ranges
func isEmojiRange2(r rune) bool {
	return (r >= 0x1F800 && r <= 0x1F8FF) || // Supplemental Arrows-C
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA00 && r <= 0x1FAFF) || // Symbols and Pictographs Extended-A
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) // Dingbats
}

// runeDisplayWidth returns the number of terminal columns used by r
func runeDisplayWidth(r rune) int {
	// Zero-width characters
	if isCombining(r) || isVariationSelector(r) || isZWJ(r) {
		return 0
	}
	// East Asian wide/fullwidth
	switch width.LookupRune(r).Kind() {
	case width.EastAsianFullwidth, width.EastAsianWide:
		return 2
	}
	// Common emoji are typically 2 columns
	if isEmoji(r) {
		return 2
	}
	return 1
}

func (e *realTimeEditor) runeWidth(r rune) int { return runeDisplayWidth(r) }

func (e *realTimeEditor) colsBetween(from, to int) int {
	if from < 0 {
		from = 0
	}
	if to < 0 {
		to = 0
	}
	if from > to {
		from, to = to, from
	}
	cols := 0
	for i := from; i < to && i < len(*e.inputRunes); i++ {
		cols += e.runeWidth((*e.inputRunes)[i])
	}
	return cols
}

func (e *realTimeEditor) moveLeft(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dD", cols)
}

func (e *realTimeEditor) moveRight(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dC", cols)
}

func (e *realTimeEditor) printTailAndReposition(from int, clearedCols int) {
	tailCols := 0
	if from < len(*e.inputRunes) {
		tail := string((*e.inputRunes)[from:])
		e.ui.write("%s", tail)
		for _, rr := range (*e.inputRunes)[from:] {
			tailCols += e.runeWidth(rr)
		}
	}
	if clearedCols > 0 {
		e.ui.write("%s", strings.Repeat(" ", clearedCols))
	}
	e.moveLeft(tailCols + clearedCols)
}

func (e *realTimeEditor) moveWordLeft() {
	if *e.cursor == 0 {
		return
	}
	i := *e.cursor - 1
	for i >= 0 && unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	for i >= 0 && !unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	newPos := i + 1
	delta := e.colsBetween(newPos, *e.cursor)
	e.moveLeft(delta)
	*e.cursor = newPos
}

func (e *realTimeEditor) moveWordRight() {
	n := len(*e.inputRunes)
	if *e.cursor >= n {
		return
	}
	i := *e.cursor
	for i < n && !unicode.IsSpace((*e.inputRunes)[i]) {
		i++
	}
	for i < n && unicode.IsSpace((*e.inputRunes)[i]) {
		i++
	}
	delta := e.colsBetween(*e.cursor, i)
	e.moveRight(delta)
	*e.cursor = i
}

// deleteWordLeft deletes the word before the cursor and updates the display
func (e *realTimeEditor) deleteWordLeft() {
	if *e.cursor == 0 {
		return
	}
	// Find new cursor position at the beginning of previous word
	i := *e.cursor - 1
	for i >= 0 && unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	for i >= 0 && !unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	newPos := i + 1
	// Compute columns to move left and columns to clear
	moveCols := e.colsBetween(newPos, *e.cursor)
	clearedCols := 0
	for j := newPos; j < *e.cursor; j++ {
		clearedCols += e.runeWidth((*e.inputRunes)[j])
	}
	// Move cursor left to newPos
	e.moveLeft(moveCols)
	// Delete runes in [newPos, cursor)
	*e.inputRunes = append((*e.inputRunes)[:newPos], (*e.inputRunes)[*e.cursor:]...)
	*e.cursor = newPos
	// Redraw tail and clear leftover cells
	e.printTailAndReposition(*e.cursor, clearedCols)
}

//nolint:revive // Input character handling inherently requires multiple cases
func (h *KeyHandler) handleInputChar(input *strings.Builder, char rune) (done bool, canceled bool) {
	switch char {
	case '\n', '\r':
		if input.Len() > 0 {
			h.ui.write("\r\n")
			return true, false
		}
		h.ui.write(" %s(required)%s", h.ui.colors.BrightRed, h.ui.colors.Reset)
		return false, false
	case '\b', 127:
		if input.Len() == 0 {
			return false, false
		}
		runes := []rune(input.String())
		if len(runes) == 0 {
			return false, false
		}
		// Identify start of previous grapheme-like cluster
		end := len(runes)
		start := end - 1
		for start >= 0 && (isCombining(runes[start]) || isVariationSelector(runes[start])) {
			start--
		}
		if start >= 0 && isRegionalIndicator(runes[start]) {
			if start > 0 && isRegionalIndicator(runes[start-1]) {
				start--
			}
		}
		for {
			if start > 0 && isZWJ(runes[start-1]) {
				start -= 2
				for start >= 0 && (isCombining(runes[start]) || isVariationSelector(runes[start])) {
					start--
				}
				continue
			}
			break
		}
		if start < 0 {
			start = 0
		}
		// Calculate total columns to clear
		cols := 0
		for i := start; i < end; i++ {
			cols += runeDisplayWidth(runes[i])
		}
		// Update input
		input.Reset()
		input.WriteString(string(runes[:start]))
		// Clear terminal cells
		for i := 0; i < cols; i++ {
			h.ui.write("\b \b")
		}
		return false, false
	case 3: // Ctrl+C
		h.ui.write("\r\n%sOperation canceled%s\r\n", h.ui.colors.BrightRed, h.ui.colors.Reset)
		return true, true
	default:
		// Accept all printable characters including multibyte
		if unicode.IsPrint(char) {
			input.WriteRune(char)
			h.ui.write("%s", string(char))
		}
		return false, false
	}
}

// getLineInput provides fallback line-based input when raw mode is not available
func (h *KeyHandler) getLineInput() (string, bool) {
	reader := bufio.NewReader(h.ui.stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", true
		}
		line = strings.TrimSpace(line)
		if line != "" {
			return line, false
		}
		h.ui.write("%s(required)%s ",
			h.ui.colors.BrightRed,
			h.ui.colors.Reset)
	}
}

// NewUI creates a new UI with the provided git client and loads keybindings from config
func NewUI(gitClient git.StatusInfoReader, router ...CommandRouter) *UI {
	colors := NewANSIColors()

	renderer := &Renderer{
		writer: os.Stdout,
		colors: colors,
	}
	renderer.updateSize()

	state := &UIState{
		selected:     0,
		input:        "",
		filtered:     []CommandInfo{},
		context:      ContextGlobal, // Start in global context
		contextStack: []Context{},
	}

	// Load config and create resolver
	var (
		cfg           *config.Config
		configManager *config.Manager
	)
	if ops, ok := gitClient.(git.ConfigOps); ok {
		configManager = config.NewConfigManager(ops)
		if err := configManager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v. Using defaults.\n", err)
		}
		cfg = configManager.GetConfig()
	} else {
		// Fallback to empty config (built-in defaults and profiles will be used)
		cfg = &config.Config{}
	}

	// Create KeyBinding resolver and register built-in profiles
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	contextManager := NewContextManager(resolver)

	// Determine which profile to use (default to "default" profile)
	profile := ProfileDefault
	if cfg.Interactive.Profile != "" {
		switch Profile(cfg.Interactive.Profile) {
		case ProfileEmacs, ProfileVi, ProfileReadline:
			profile = Profile(cfg.Interactive.Profile)
		default:
			fmt.Fprintf(os.Stderr, "Warning: Unknown profile '%s', using default\n", cfg.Interactive.Profile)
		}
	}

	// Resolve contextual keybindings for all contexts
	contextualMap, err := resolver.ResolveContextual(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to resolve keybindings: %v. Using defaults.\n", err)
		// Fallback to legacy defaults
		keyMap := DefaultKeyBindingMap()
		contextualMap = &ContextualKeyBindingMap{
			Profile:  ProfileDefault,
			Platform: DetectPlatform(),
			Terminal: DetectTerminal(),
			Contexts: map[Context]*KeyBindingMap{
				ContextGlobal:            keyMap,
				ContextInput:             keyMap,
				ContextResults:           keyMap,
				ContextSearch:            keyMap,
				ContextWorkflowView:      keyMap,
				ContextWorkflowSelection: keyMap,
			},
		}
	}

	ui := &UI{
		stdin:       os.Stdin,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
		term:        termio.DefaultTerminal{},
		renderer:    renderer,
		state:       state,
		colors:      colors,
		gitClient:   gitClient,
		gitStatus:   getGitStatus(gitClient),
		contextMgr:  contextManager,
		profile:     profile,
		workflowMgr: NewWorkflowManager(),
		configMgr:   configManager,
		config:      cfg,
	}
	state.onContextChange = func(_ Context, newCtx Context) {
		contextManager.SetContext(newCtx)
	}
	ui.bootstrapConfigWorkflows()
	state.SetWorkflowListIndex(0, len(ui.workflowMgr.ListWorkflows()))
	ui.updateWorkflowPointer()

	ui.handler = &KeyHandler{
		ui:            ui,
		contextualMap: contextualMap,
	}

	// Set up workflow executor if router is provided
	if len(router) > 0 && router[0] != nil {
		ui.workflowEx = NewWorkflowExecutor(router[0], ui)
	}

	return ui
}

// ToggleWorkflowView toggles between normal command view and workflow view
func (ui *UI) ToggleWorkflowView() {
	ui.state.showWorkflow = !ui.state.showWorkflow
	if ui.state.showWorkflow {
		ui.state.EnterContext(ContextWorkflowView)
		ui.state.DeactivateWorkflowSelection()
		ui.clearPendingWorkflowSelection()
		ui.ensureWorkflowListSelection()
	} else {
		ui.state.ExitContext()
		ui.state.SetWorkflowListIndex(0, len(ui.listWorkflows()))
		ui.updateWorkflowPointer()
	}
}

// AddToWorkflow adds a command to the workflow
func (ui *UI) AddToWorkflow(command string, args []string, description string) int {
	id, _ := ui.AddToWorkflowByID(ui.workflowMgr.GetActiveID(), command, args, description)
	return id
}

// AddToWorkflowByID adds a step to the specified workflow.
func (ui *UI) AddToWorkflowByID(workflowID int, command string, args []string, description string) (int, error) {
	if ui.workflowMgr == nil {
		return 0, fmt.Errorf("workflow manager not initialized")
	}

	if workflowID == 0 {
		workflowID = ui.workflowMgr.GetActiveID()
	}

	id, err := ui.workflowMgr.AddStep(workflowID, command, args, description)
	if err != nil {
		return 0, err
	}

	ui.ensureWorkflowListSelection()
	ui.updateWorkflowPointer()
	return id, nil
}

// ApplyContextualKeybindings updates the active keybinding map, satisfying keybindings.ContextualMapApplier.
func (ui *UI) ApplyContextualKeybindings(contextual *ContextualKeyBindingMap) {
	if ui == nil || ui.handler == nil || contextual == nil {
		return
	}
	ui.handler.contextualMap = contextual
}

func (ui *UI) resetToSearchMode() bool {
	if ui == nil || ui.state == nil {
		return false
	}

	state := ui.state
	active := state.HasInput() || state.showWorkflow || state.workflowSelectionActive || len(state.contextStack) > 0 || state.GetCurrentContext() != ContextGlobal
	state.ClearInput()
	state.selected = 0
	state.contextStack = nil
	state.SetContext(ContextGlobal)
	state.showWorkflow = false
	state.DeactivateWorkflowSelection()
	ui.clearPendingWorkflowSelection()
	return active
}

// ResetToSearchMode clears the interactive search UI back to its default state.
func (ui *UI) ResetToSearchMode() bool {
	return ui.resetToSearchMode()
}

func (ui *UI) readPlaceholderInput() (string, bool) {
	if ui == nil || ui.handler == nil {
		return "", true
	}
	return ui.handler.getRealTimeInput()
}

// ClearWorkflow removes all steps from the workflow
func (ui *UI) ClearWorkflow() {
	_ = ui.ClearWorkflowByID(ui.workflowMgr.GetActiveID())
}

// ClearWorkflowByID clears the workflow with the provided ID.
func (ui *UI) ClearWorkflowByID(workflowID int) error {
	if ui.workflowMgr == nil {
		return fmt.Errorf("workflow manager not initialized")
	}
	if !ui.workflowMgr.ClearWorkflow(workflowID) {
		return fmt.Errorf("workflow %d not found", workflowID)
	}
	ui.updateWorkflowPointer()
	return nil
}

// ExecuteWorkflow executes the current workflow
func (ui *UI) ExecuteWorkflow() error {
	return ui.ExecuteWorkflowByID(ui.workflowMgr.GetActiveID())
}

// ExecuteWorkflowByID executes the workflow identified by workflowID.
func (ui *UI) ExecuteWorkflowByID(workflowID int) error {
	if ui.workflowEx == nil {
		return fmt.Errorf("workflow executor not initialized")
	}

	workflow, exists := ui.workflowMgr.GetWorkflow(workflowID)
	if !exists {
		return fmt.Errorf("workflow %d not found", workflowID)
	}

	if workflow.IsEmpty() {
		return fmt.Errorf("workflow is empty")
	}

	return ui.workflowEx.Execute(workflow)
}

// listWorkflows provides a snapshot of available workflows for rendering.
func (ui *UI) listWorkflows() []WorkflowSummary {
	if ui == nil || ui.workflowMgr == nil {
		return nil
	}
	return ui.workflowMgr.ListWorkflows()
}

// ensureWorkflowListSelection keeps the management view index aligned with the active workflow.
func (ui *UI) ensureWorkflowListSelection() {
	if ui == nil {
		return
	}
	summaries := ui.listWorkflows()
	if len(summaries) == 0 {
		return
	}

	activeID := ui.workflowMgr.GetActiveID()
	index := 0
	for i, summary := range summaries {
		if summary.ID == activeID {
			index = i
			break
		}
	}
	ui.state.SetWorkflowListIndex(index, len(summaries))
	ui.updateWorkflowPointer()
}

// clearPendingWorkflowSelection releases transient command data captured by the selection overlay.
func (ui *UI) clearPendingWorkflowSelection() {
	ui.pendingWorkflowTemplate = ""
	ui.pendingWorkflowCommand = ""
	ui.pendingWorkflowArgs = nil
}

// workflowStepsByID returns a copy of workflow steps for rendering purposes.
func (ui *UI) workflowStepsByID(workflowID int) []WorkflowStep {
	if ui == nil || ui.workflowMgr == nil {
		return nil
	}
	workflow, exists := ui.workflowMgr.GetWorkflow(workflowID)
	if !exists || workflow == nil {
		return nil
	}
	return workflow.GetSteps()
}

// hasWorkflowSteps indicates whether any workflow currently contains steps.
func (ui *UI) hasWorkflowSteps() bool {
	for _, summary := range ui.listWorkflows() {
		if summary.StepCount > 0 {
			return true
		}
	}
	return false
}

func (ui *UI) updateWorkflowPointer() {
	if ui == nil || ui.workflowMgr == nil {
		ui.workflow = nil
		return
	}
	workflow, _ := ui.workflowMgr.GetActiveWorkflow()
	ui.workflow = workflow
}

func (ui *UI) bootstrapConfigWorkflows() {
	if ui == nil || ui.workflowMgr == nil || ui.config == nil {
		return
	}

	for _, wf := range ui.config.Workflows {
		if len(wf.Steps) == 0 {
			continue
		}
		if _, err := ui.workflowMgr.CreateReadOnlyWorkflow(wf.Name, wf.Steps); err != nil {
			_, _ = fmt.Fprintf(ui.stderr, "Warning: failed to load workflow %q from config: %v\n", wf.Name, err)
		}
	}
	ui.ensureWorkflowListSelection()
	ui.updateWorkflowPointer()
}

// updateSize updates the terminal dimensions
func (r *Renderer) updateSize() {
	w, h := uiutil.Dimensions(r.writer, 80, 24)
	r.width, r.height = w, h
}

var commands = buildInteractiveCommands()

// Run executes the incremental search interactive UI with the provided custom git client,
// and returns the selected command as []string (or nil if nothing is selected).
func Run(gitClient git.StatusInfoReader) []string {
	ui := NewUI(gitClient)
	return ui.Run()
}

// writeError writes an error message to stderr
func (ui *UI) writeError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stderr, format+"\n", a...)
}

// write writes a message to stdout
func (ui *UI) write(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stdout, format, a...)
}

// writeColor writes a colored message to stdout
func (ui *UI) writeColor(text string) {
	_, _ = fmt.Fprint(ui.stdout, text)
}

// writeln writes a message with newline to stdout
func (ui *UI) writeln(format string, a ...interface{}) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(ui.stdout, "\r\x1b[K")
	_, _ = fmt.Fprintf(ui.stdout, format+"\r\n", a...)
}

func (ui *UI) notifySoftCancel() {
	ui.softCancelFlash.Store(true)
}

func (ui *UI) consumeSoftCancelFlash() bool {
	return ui.softCancelFlash.Swap(false)
}

// clearScreen clears the entire screen and hides cursor
func clearScreen(w io.Writer) {
	uiutil.ClearScreen(w)
	uiutil.HideCursor(w)
}

// showCursor shows the terminal cursor
func showCursor(w io.Writer) {
	uiutil.ShowCursor(w)
}

// ellipsis truncates string and adds ellipsis if it exceeds maxLen (ASCII only)
func ellipsis(s string, maxLen int) string {
	return uiutil.Ellipsis(s, maxLen)
}

// Render displays the command list with proper terminal handling
func (r *Renderer) Render(ui *UI, state *UIState) {
	clearScreen(r.writer)
	// Disable line wrapping during rendering, restore at end
	uiutil.DisableWrap(r.writer)
	var restoreCursor func()
	defer func() {
		uiutil.EnableWrap(r.writer)
		if restoreCursor != nil {
			restoreCursor()
		}
		showCursor(r.writer)
	}()

	// Update terminal size
	r.updateSize()

	// Render each section
	r.renderHeader(ui)
	r.renderSoftCancelFlash(ui)

	r.renderWorkflowStatusIfNeeded(ui, state)

	if state.showWorkflow {
		r.renderWorkflowView(ui, state)
		return
	}

	restoreCursor = r.renderSearchArea(ui, state)
	if state.workflowSelectionActive {
		r.renderWorkflowSelection(ui, state)
	}
}

func (r *Renderer) renderWorkflowStatusIfNeeded(ui *UI, state *UIState) {
	if !ui.hasWorkflowSteps() || state.showWorkflow || state.workflowSelectionActive || state.input != "" {
		return
	}
	r.renderWorkflowStatus(ui)
}

func (r *Renderer) renderSearchArea(ui *UI, state *UIState) func() {
	r.renderSearchPrompt(ui, state)
	restoreCursor := r.saveCursorAtSearchPrompt(state)

	switch {
	case state.input == "":
		r.renderInitialSearchState(ui)
	case len(state.filtered) == 0:
		r.renderNoMatches(ui, state)
	default:
		r.renderCommandList(ui, state)
	}

	return restoreCursor
}

func (r *Renderer) renderInitialSearchState(ui *UI) {
	if !ui.hasWorkflowSteps() {
		r.renderEmptyState(ui)
	} else {
		r.renderEmptyStateWithWorkflow(ui)
	}
	r.writeEmptyLine()
	r.renderSearchKeybinds(ui)
}

func (r *Renderer) renderSoftCancelFlash(ui *UI) {
	if !ui.consumeSoftCancelFlash() {
		return
	}
	alert := fmt.Sprintf("%s⚠️  Operation canceled%s", r.colors.BrightRed+r.colors.Bold, r.colors.Reset)
	r.writeColorln(ui, alert)
	r.writeColorln(ui, "")
}

// renderHeader renders the title, git status, and navigation subtitle
func (r *Renderer) renderHeader(ui *UI) {
	// Modern header with title
	title := fmt.Sprintf("%s%s🚀 ggc Interactive Mode%s",
		r.colors.BrightCyan+r.colors.Bold,
		r.colors.Reset,
		r.colors.Reset)
	r.writeColorln(ui, title)

	// Git status information
	if ui.gitStatus != nil {
		r.renderGitStatus(ui, ui.gitStatus)
	}

}

// renderSearchPrompt renders the search input with cursor
func (r *Renderer) renderSearchPrompt(ui *UI, state *UIState) {
	inputWithCursor := r.formatInputWithCursor(state)

	searchPrompt := fmt.Sprintf("%s┌─ %sSearch:%s %s",
		r.colors.BrightBlue,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.Reset,
		inputWithCursor)
	r.writeColorln(ui, searchPrompt)

	// Results separator
	if state.input != "" {
		separator := fmt.Sprintf("%s└─ %sResults:%s",
			r.colors.BrightBlue,
			r.colors.BrightMagenta+r.colors.Bold,
			r.colors.Reset)
		r.writeColorln(ui, separator)
	}
	r.writeEmptyLine()
}

func (r *Renderer) saveCursorAtSearchPrompt(state *UIState) func() {
	linesUp := 2
	if state.input != "" {
		linesUp++
	}
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dA", linesUp)
	const prefix = "┌─ Search: "
	// Compute display width (columns) of the prefix using runeDisplayWidth
	prefixCols := 0
	for _, pr := range prefix {
		prefixCols += runeDisplayWidth(pr)
	}
	// Compute display width up to the logical cursor position
	runes := []rune(state.input)
	cursorPos := state.cursorPos
	if cursorPos > len(runes) {
		cursorPos = len(runes)
	}
	cursorWidth := 0
	for _, rr := range runes[:cursorPos] {
		cursorWidth += runeDisplayWidth(rr)
	}
	column := prefixCols + cursorWidth + 1
	if column < 1 {
		column = 1
	}
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dG", column)
	_, _ = fmt.Fprint(r.writer, "\x1b[s")
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dB", linesUp)
	return func() {
		_, _ = fmt.Fprint(r.writer, "\x1b[u")
	}
}

// formatInputWithCursor formats the input string with cursor position
func (r *Renderer) formatInputWithCursor(state *UIState) string {
	if state.input == "" {
		return fmt.Sprintf("%s█%s", r.colors.BrightWhite+r.colors.Bold, r.colors.Reset)
	}

	inputRunes := []rune(state.input)
	beforeCursor := string(inputRunes[:state.cursorPos])
	afterCursor := string(inputRunes[state.cursorPos:])
	cursor := "│"
	if state.cursorPos >= utf8.RuneCountInString(state.input) {
		cursor = "█"
	}

	return fmt.Sprintf("%s%s%s%s%s%s%s",
		r.colors.BrightYellow,
		beforeCursor,
		r.colors.BrightWhite+r.colors.Bold,
		cursor,
		r.colors.Reset+r.colors.BrightYellow,
		afterCursor,
		r.colors.Reset)
}

// renderEmptyState renders the empty input state
func (r *Renderer) renderEmptyState(ui *UI) {
	r.writeColorln(ui, fmt.Sprintf("%s💭 %sStart typing to search commands...%s",
		r.colors.BrightBlue, r.colors.BrightBlack, r.colors.Reset))
}

func (r *Renderer) buildSearchKeybindEntries(ui *UI) []keybindHelpEntry {
	entries := []keybindHelpEntry{
		{key: "←/→", desc: "Move cursor"},
		{key: "Ctrl+←/→", desc: "Move by word"},
		{key: "Option+←/→", desc: "Move by word (macOS)"},
	}
	// Future: extend this helper for additional contexts such as workflow views.

	var km *KeyBindingMap
	if ui != nil && ui.handler != nil {
		km = ui.handler.GetCurrentKeyMap()
	}
	if km == nil {
		km = DefaultKeyBindingMap()
	}

	defaultMap := DefaultKeyBindingMap()

	appendDynamic := func(primary []KeyStroke, fallback []KeyStroke, desc string) {
		keys := primary
		if len(keys) == 0 {
			keys = fallback
		}
		if len(keys) == 0 {
			return
		}
		formatted := FormatKeyStrokesForDisplay(keys)
		if formatted == "" || formatted == "none" {
			return
		}
		entries = append(entries, keybindHelpEntry{key: formatted, desc: desc})
	}

	appendDynamic(km.ClearLine, defaultMap.ClearLine, "Clear all input")
	appendDynamic(km.DeleteWord, defaultMap.DeleteWord, "Delete word")
	appendDynamic(km.DeleteToEnd, defaultMap.DeleteToEnd, "Delete to end")
	appendDynamic(km.MoveToBeginning, defaultMap.MoveToBeginning, "Move to beginning")
	appendDynamic(km.MoveToEnd, defaultMap.MoveToEnd, "Move to end")

	entries = append(entries, keybindHelpEntry{key: "Backspace", desc: "Delete character"})
	entries = append(entries, keybindHelpEntry{key: "Enter", desc: "Execute selected command"})

	appendDynamic(km.AddToWorkflow, defaultMap.AddToWorkflow, "Add to workflow")
	appendDynamic(km.ToggleWorkflowView, defaultMap.ToggleWorkflowView, "Toggle workflow view")

	entries = append(entries, keybindHelpEntry{key: "Ctrl+c", desc: "Quit"})

	return entries
}

func (r *Renderer) renderKeybindEntries(ui *UI, entries []keybindHelpEntry) {
	if len(entries) == 0 {
		return
	}

	r.writeColorln(ui, fmt.Sprintf("%s⌨️  %sAvailable keybinds:%s",
		r.colors.BrightBlue, r.colors.BrightWhite+r.colors.Bold, r.colors.Reset))

	for _, entry := range entries {
		r.writeColorln(ui, fmt.Sprintf("   %s%s%s  %s%s%s",
			r.colors.BrightGreen+r.colors.Bold,
			entry.key,
			r.colors.Reset,
			r.colors.BrightBlack,
			entry.desc,
			r.colors.Reset))
	}
}

// renderNoMatches renders the no matches found state with keybind help
func (r *Renderer) renderNoMatches(ui *UI, state *UIState) {
	// No matches message
	r.writeColorln(ui, fmt.Sprintf("%s🔍 %sNo commands found for '%s%s%s'%s",
		r.colors.BrightYellow,
		r.colors.BrightWhite,
		r.colors.BrightYellow+r.colors.Bold,
		state.input,
		r.colors.Reset+r.colors.BrightWhite,
		r.colors.Reset))

	r.writeEmptyLine()
	r.renderKeybindEntries(ui, r.buildSearchKeybindEntries(ui))
}

// renderSearchKeybinds renders keybinds available in search UI
func (r *Renderer) renderSearchKeybinds(ui *UI) {
	r.renderKeybindEntries(ui, r.buildSearchKeybindEntries(ui))
}

// renderWorkflowKeybinds renders keybinds available in workflow UI
func (r *Renderer) renderWorkflowKeybinds(ui *UI) {
	// Workflow-specific keybinds
	defaultMap := DefaultKeyBindingMap()
	km := defaultMap
	if ui != nil && ui.handler != nil {
		if current := ui.handler.GetCurrentKeyMap(); current != nil {
			km = current
		}
	}

	format := func(keys []KeyStroke, fallback []KeyStroke, placeholder string) string {
		selected := keys
		if len(selected) == 0 {
			selected = fallback
		}
		if len(selected) == 0 {
			return placeholder
		}
		if formatted := FormatKeyStrokesForDisplay(selected); formatted != "" {
			return formatted
		}
		return placeholder
	}

	keybinds := []struct{ key, desc string }{
		{"↑/↓", "Navigate workflows"},
		{"Enter", "Execute selected workflow"},
		{format(km.WorkflowCreate, defaultMap.WorkflowCreate, "Ctrl+n"), "Create new workflow"},
		{format(km.WorkflowDelete, defaultMap.WorkflowDelete, "d"), "Delete selected workflow"},
		{format(km.WorkflowCancel, defaultMap.WorkflowCancel, "Esc"), "Back to search view"},
		{format(km.ToggleWorkflowView, defaultMap.ToggleWorkflowView, "Ctrl+t"), "Back to search view"},
		{"Ctrl+c", "Quit"},
	}

	r.writeColorln(ui, fmt.Sprintf("%s⌨️  %sAvailable keybinds:%s",
		r.colors.BrightBlue, r.colors.BrightWhite+r.colors.Bold, r.colors.Reset))

	for _, kb := range keybinds {
		r.writeColorln(ui, fmt.Sprintf("   %s%s%s  %s%s%s",
			r.colors.BrightGreen+r.colors.Bold,
			kb.key,
			r.colors.Reset,
			r.colors.BrightBlack,
			kb.desc,
			r.colors.Reset))
	}
}

// renderCommandList renders the filtered command list
func (r *Renderer) renderCommandList(ui *UI, state *UIState) {
	// Clamp selection index to valid range
	if state.selected >= len(state.filtered) {
		state.selected = len(state.filtered) - 1
	}
	if state.selected < 0 {
		state.selected = 0
	}

	// Calculate maximum command length for consistent alignment
	maxCmdLen := r.calculateMaxCommandLength(state.filtered)

	for i, cmd := range state.filtered {
		r.renderCommandItem(ui, cmd, i, state.selected, maxCmdLen)
	}
}

// renderCommandItem renders a single command item
func (r *Renderer) renderCommandItem(ui *UI, cmd CommandInfo, index, selected, maxCmdLen int) {
	desc := cmd.Description
	if desc == "" {
		desc = "No description"
	}

	// Calculate padding for consistent command alignment
	paddingLen := maxCmdLen - len(cmd.Command)
	if paddingLen < 0 {
		paddingLen = 0
	}
	padding := strings.Repeat(" ", paddingLen)

	// Calculate available width for description
	usedWidth := 4 + len(cmd.Command) + len(padding) + 3 // prefix + command + padding + separator
	availableDescWidth := r.width - usedWidth
	if availableDescWidth < 10 {
		availableDescWidth = 10
	}

	// Truncate description if needed
	trimmedDesc := ellipsis(desc, availableDescWidth)

	if index == selected {
		// Selected item with modern highlighting
		selectedLine := fmt.Sprintf("%s▶ %s%s%s%s %s│%s %s%s%s",
			r.colors.BrightCyan+r.colors.Bold,
			r.colors.BrightWhite+r.colors.Bold+r.colors.Reverse,
			" "+cmd.Command+" ",
			r.colors.Reset,
			padding,
			r.colors.BrightBlue,
			r.colors.Reset,
			r.colors.BrightWhite,
			trimmedDesc,
			r.colors.Reset)
		r.writeColorln(ui, selectedLine)
	} else {
		// Regular item with improved styling
		regularLine := fmt.Sprintf("  %s%s%s%s %s│%s %s%s%s",
			r.colors.BrightGreen+r.colors.Bold,
			cmd.Command,
			r.colors.Reset,
			padding,
			r.colors.BrightBlack,
			r.colors.Reset,
			r.colors.BrightBlack,
			trimmedDesc,
			r.colors.Reset)
		r.writeColorln(ui, regularLine)
	}
}

// writeColorln writes a colored line to the terminal
func (r *Renderer) writeColorln(_ *UI, text string) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(r.writer, "\r\x1b[K")
	_, _ = fmt.Fprint(r.writer, text+"\r\n")
}

// renderEmptyStateWithWorkflow renders the empty state with workflow info
func (r *Renderer) renderEmptyStateWithWorkflow(ui *UI) {
	activeID := ui.workflowMgr.GetActiveID()
	summaries := ui.listWorkflows()

	var activeSteps int
	totalSteps := 0
	for _, summary := range summaries {
		totalSteps += summary.StepCount
		if summary.ID == activeID {
			activeSteps = summary.StepCount
		}
	}

	r.writeColorln(ui, fmt.Sprintf("%s📝 Workflows ready. Press Tab to choose where the next step goes.%s",
		r.colors.BrightBlue,
		r.colors.Reset))
	r.writeColorln(ui, fmt.Sprintf("   %sActive workflow #%d:%s %s%d step(s)%s | %sTotal:%s %s%d%s",
		r.colors.BrightBlack,
		activeID,
		r.colors.Reset,
		r.colors.BrightGreen,
		activeSteps,
		r.colors.Reset,
		r.colors.BrightBlack,
		r.colors.Reset,
		r.colors.BrightGreen,
		totalSteps,
		r.colors.Reset))
}

// renderWorkflowStatus renders workflow information at the top of the UI
func (r *Renderer) renderWorkflowStatus(ui *UI) {
	summaries := ui.listWorkflows()
	if len(summaries) == 0 {
		return
	}

	activeSummary, totalSteps, otherWithSteps := summarizeWorkflowStatus(summaries)

	if totalSteps == 0 {
		return
	}

	desc := r.describeActiveWorkflow(activeSummary)

	statusText := fmt.Sprintf("%s📋 %s%s%s %s(%d step(s))%s",
		r.colors.BrightYellow+r.colors.Bold,
		desc.TagSegment,
		r.colors.BrightWhite+r.colors.Bold,
		desc.Name,
		r.colors.BrightBlack,
		desc.Steps,
		r.colors.Reset)

	steps := ui.workflowStepsByID(desc.ID)
	stepTexts := make([]string, 0, min(3, len(steps)))
	for i, step := range steps[:min(3, len(steps))] {
		stepText := fmt.Sprintf("%s%d.%s %s%s%s",
			r.colors.BrightBlue+r.colors.Bold,
			i+1,
			r.colors.Reset,
			r.colors.BrightGreen,
			step.Description,
			r.colors.Reset)
		stepTexts = append(stepTexts, stepText)
	}

	if len(steps) > 3 {
		stepTexts = append(stepTexts, fmt.Sprintf("%s...+%d more%s",
			r.colors.BrightBlack, len(steps)-3, r.colors.Reset))
	}

	if otherWithSteps > 0 {
		stepTexts = append(stepTexts, fmt.Sprintf("%s+%d other workflow(s)%s",
			r.colors.BrightBlack, otherWithSteps, r.colors.Reset))
	}

	r.writeColorln(ui, statusText+" "+strings.Join(stepTexts, " → "))
	r.writeColorln(ui, fmt.Sprintf("%sTotal steps across workflows:%s %s%d%s",
		r.colors.BrightBlack,
		r.colors.Reset,
		r.colors.BrightGreen,
		totalSteps,
		r.colors.Reset))
	r.writeColorln(ui, "")
}

func summarizeWorkflowStatus(summaries []WorkflowSummary) (*WorkflowSummary, int, int) {
	var activeSummary *WorkflowSummary
	totalSteps := 0
	otherWithSteps := 0

	for i := range summaries {
		summary := &summaries[i]
		totalSteps += summary.StepCount
		if summary.IsActive {
			activeSummary = summary
		} else if summary.StepCount > 0 {
			otherWithSteps++
		}
	}

	return activeSummary, totalSteps, otherWithSteps
}

type activeWorkflowDescriptor struct {
	ID         int
	Steps      int
	Name       string
	TagSegment string
}

func (r *Renderer) describeActiveWorkflow(activeSummary *WorkflowSummary) activeWorkflowDescriptor {
	desc := activeWorkflowDescriptor{
		Name: "Active workflow",
	}

	if activeSummary == nil {
		return desc
	}

	desc.ID = activeSummary.ID
	desc.Steps = activeSummary.StepCount
	if trimmed := strings.TrimSpace(activeSummary.Name); trimmed != "" {
		desc.Name = trimmed
	} else {
		desc.Name = fmt.Sprintf("Workflow #%d", activeSummary.ID)
	}

	tagText := "[Dynamic]"
	tagColor := r.colors.BrightMagenta + r.colors.Bold
	if activeSummary.Source == WorkflowSourceConfig {
		tagText = "[Config]"
		tagColor = r.colors.BrightYellow + r.colors.Bold
	}

	if tagText != "" {
		desc.TagSegment = fmt.Sprintf("%s%s%s ", tagColor, tagText, r.colors.Reset)
	}

	return desc
}

// renderWorkflowView renders the detailed workflow view
func (r *Renderer) renderWorkflowView(ui *UI, state *UIState) {
	summaries := ui.listWorkflows()

	// Detailed workflow header
	r.writeColorln(ui, fmt.Sprintf("%s📋 Multiple Workflows%s",
		r.colors.BrightYellow+r.colors.Bold,
		r.colors.Reset))
	r.writeColorln(ui, "")

	if len(summaries) == 0 {
		r.writeColorln(ui, fmt.Sprintf("%s  No steps in workflow%s",
			r.colors.BrightBlack,
			r.colors.Reset))
		r.writeColorln(ui, "")

		// Render workflow keybinds even for empty workflow
		r.renderWorkflowKeybinds(ui)
		return
	}

	for i, summary := range summaries {
		isSelected := i == state.workflowListIndex
		prefix := "    "
		if isSelected {
			prefix = fmt.Sprintf("%s▶ %s", r.colors.BrightCyan+r.colors.Bold, r.colors.Reset)
		}

		displayName := strings.TrimSpace(summary.Name)
		if displayName == "" {
			displayName = fmt.Sprintf("Workflow #%d", summary.ID)
		}

		tagText := "[Dynamic]"
		tagColor := r.colors.BrightMagenta + r.colors.Bold
		if summary.Source == WorkflowSourceConfig {
			tagText = "[Config]"
			tagColor = r.colors.BrightYellow + r.colors.Bold
		}

		metaLine := fmt.Sprintf("%s%s%s %s%s%s %s(%d step(s))%s",
			tagColor,
			tagText,
			r.colors.Reset,
			r.colors.BrightGreen+r.colors.Bold,
			displayName,
			r.colors.Reset,
			r.colors.BrightBlack,
			summary.StepCount,
			r.colors.Reset)

		if summary.IsActive {
			metaLine += fmt.Sprintf(" %s[active]%s",
				r.colors.BrightBlue,
				r.colors.Reset)
		}

		r.writeColorln(ui, prefix+metaLine)

		steps := ui.workflowStepsByID(summary.ID)
		if len(steps) == 0 {
			r.writeColorln(ui, fmt.Sprintf("%s      (empty)%s",
				r.colors.BrightBlack,
				r.colors.Reset))
		} else {
			for idx, step := range steps {
				stepLine := fmt.Sprintf("      %s%d.%s %s%s%s",
					r.colors.BrightBlue+r.colors.Bold,
					idx+1,
					r.colors.Reset,
					r.colors.BrightGreen+r.colors.Bold,
					step.Description,
					r.colors.Reset)
				r.writeColorln(ui, stepLine)
			}
		}

		r.writeColorln(ui, "")
	}

	// Render workflow keybinds
	r.renderWorkflowKeybinds(ui)
}

// renderWorkflowSelection shows the overlay for choosing a workflow target.
func (r *Renderer) renderWorkflowSelection(ui *UI, state *UIState) {
	summaries := ui.listWorkflows()
	totalOptions := len(summaries) + 1 // include "create new"
	if totalOptions == 0 {
		return
	}

	r.writeColorln(ui, fmt.Sprintf("%s┌─ Select Workflow ──────────────────────────┐%s",
		r.colors.BrightYellow,
		r.colors.Reset))

	for idx, summary := range summaries {
		selected := state.workflowSelectionIndex == idx
		prefix := "│   "
		if selected {
			prefix = fmt.Sprintf("│ %s▶%s ", r.colors.BrightCyan+r.colors.Bold, r.colors.Reset)
		}

		displayName := strings.TrimSpace(summary.Name)
		if displayName == "" {
			displayName = fmt.Sprintf("Workflow #%d", summary.ID)
		}

		tagText := "[Dynamic]"
		tagColor := r.colors.BrightMagenta + r.colors.Bold
		if summary.Source == WorkflowSourceConfig {
			tagText = "[Config]"
			tagColor = r.colors.BrightYellow + r.colors.Bold
		}
		tagSegment := fmt.Sprintf("%s%s%s", tagColor, tagText, r.colors.Reset)

		line := fmt.Sprintf("%s %s%s%s %s(%d step%s)%s",
			tagSegment,
			r.colors.BrightWhite+r.colors.Bold,
			displayName,
			r.colors.Reset,
			r.colors.BrightBlack,
			summary.StepCount,
			pluralize(summary.StepCount),
			r.colors.Reset)

		if summary.IsActive {
			line += fmt.Sprintf(" %s[active]%s", r.colors.BrightBlue, r.colors.Reset)
		}

		r.writeColorln(ui, fmt.Sprintf("%s%s │",
			prefix,
			line))
	}

	createIdx := len(summaries)
	selectedCreate := state.workflowSelectionIndex == createIdx
	prefix := "│   "
	if selectedCreate {
		prefix = fmt.Sprintf("│ %s▶%s ", r.colors.BrightCyan+r.colors.Bold, r.colors.Reset)
	}
	createLabel := fmt.Sprintf("%s+ Create new workflow%s",
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.Reset)
	r.writeColorln(ui, fmt.Sprintf("%s%-36s│", prefix, createLabel))

	r.writeColorln(ui, fmt.Sprintf("%s└────────────────────────────────────────────┘%s",
		r.colors.BrightYellow,
		r.colors.Reset))

	r.writeColorln(ui, fmt.Sprintf("%s   Enter%s to confirm • %sEsc%s to cancel",
		r.colors.BrightBlack,
		r.colors.Reset,
		r.colors.BrightBlack,
		r.colors.Reset))
	r.writeColorln(ui, "")
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// renderGitStatus renders the Git repository status information
func (r *Renderer) renderGitStatus(ui *UI, status *GitStatus) {
	var parts []string

	// Branch name
	branchPart := fmt.Sprintf("%s📍 %s%s%s",
		r.colors.BrightBlue,
		r.colors.BrightWhite+r.colors.Bold,
		status.Branch,
		r.colors.Reset)
	parts = append(parts, branchPart)

	// Working directory status
	if status.HasChanges {
		var statusParts []string
		if status.Modified > 0 {
			statusParts = append(statusParts, fmt.Sprintf("%d modified", status.Modified))
		}
		if status.Staged > 0 {
			statusParts = append(statusParts, fmt.Sprintf("%d staged", status.Staged))
		}

		workingPart := fmt.Sprintf("%s📝 %s%s%s",
			r.colors.BrightYellow,
			r.colors.BrightWhite+r.colors.Bold,
			strings.Join(statusParts, ", "),
			r.colors.Reset)
		parts = append(parts, workingPart)
	}

	// Remote tracking status
	if status.Ahead > 0 || status.Behind > 0 {
		var remoteParts []string
		if status.Ahead > 0 {
			remoteParts = append(remoteParts, fmt.Sprintf("↑%d", status.Ahead))
		}
		if status.Behind > 0 {
			remoteParts = append(remoteParts, fmt.Sprintf("↓%d", status.Behind))
		}

		remotePart := fmt.Sprintf("%s%s%s",
			r.colors.BrightMagenta+r.colors.Bold,
			strings.Join(remoteParts, " "),
			r.colors.Reset)
		parts = append(parts, remotePart)
	}

	// Render the status line
	statusLine := strings.Join(parts, "  ")
	r.writeColorln(ui, statusLine)
}

// writeEmptyLine writes an empty line
func (r *Renderer) writeEmptyLine() {
	_, _ = fmt.Fprint(r.writer, "\r\x1b[K\r\n")
}

// calculateMaxCommandLength calculates the maximum command length for alignment
func (r *Renderer) calculateMaxCommandLength(filtered []CommandInfo) int {
	if len(filtered) == 0 {
		return 0
	}

	maxLen := 0
	for _, cmd := range filtered {
		if len(cmd.Command) > maxLen {
			maxLen = len(cmd.Command)
		}
	}
	return maxLen
}

// setupTerminal configures terminal raw mode and returns the old state and error status
func (ui *UI) setupTerminal() (*term.State, bool) {
	var oldState *term.State
	if f, ok := ui.stdin.(*os.File); ok {
		// Check if it's a real terminal (TTY) - skip for non-TTY or when term is nil
		if ui.term == nil || !term.IsTerminal(int(f.Fd())) {
			// Not a real terminal, skip raw mode setup for debugging
			return nil, true
		}

		fd := int(f.Fd())
		var err error
		oldState, err = ui.term.MakeRaw(fd)
		if err != nil {
			ui.writeError("Failed to set terminal to raw mode: %v", err)
			return nil, false
		}
	}
	return oldState, true
}

// Run executes the interactive UI
func (ui *UI) Run() []string {
	oldState, reader, isRawMode := ui.initializeTerminal()
	if oldState == nil && isRawMode {
		return nil
	}

	// Set up terminal restoration for raw mode
	if f, ok := ui.stdin.(*os.File); ok && isRawMode {
		fd := int(f.Fd())
		defer func() {
			if err := ui.term.Restore(fd, oldState); err != nil {
				ui.writeError("failed to restore terminal state: %v", err)
			}
		}()
	}

	return ui.runMainLoop(reader, isRawMode, oldState)
}

// initializeTerminal sets up the terminal and returns the old state, reader, and raw mode status
func (ui *UI) initializeTerminal() (*term.State, *bufio.Reader, bool) {
	oldState, ok := ui.setupTerminal()
	if !ok {
		return nil, nil, false
	}

	// Check if we're in raw mode (real terminal) or not
	isRawMode := oldState != nil

	// Set up reader based on mode
	var reader *bufio.Reader
	if !isRawMode {
		reader = bufio.NewReader(ui.stdin)
	}

	return oldState, reader, isRawMode
}

// runMainLoop handles the main input loop
func (ui *UI) runMainLoop(reader *bufio.Reader, isRawMode bool, oldState *term.State) []string {
	if isRawMode {
		if ui.reader == nil {
			ui.reader = bufio.NewReader(ui.stdin)
		}
	} else {
		ui.reader = reader
	}

	for {
		ui.state.UpdateFiltered()
		ui.renderer.Render(ui, ui.state)

		r, err := ui.readNextRune(reader, isRawMode)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			continue // Skip this iteration for other errors
		}

		// Handle key input with rune
		isSingleByte := isRawMode // In raw mode, we read single bytes; in buffered mode, we read full runes
		shouldContinue, result := ui.handler.HandleKey(r, isSingleByte, oldState, reader)
		if !shouldContinue {
			return result
		}
	}
}

// readNextRune reads the next rune from input based on the mode
func (ui *UI) readNextRune(reader *bufio.Reader, isRawMode bool) (rune, error) {
	if isRawMode {
		for {
			// Read single byte directly from stdin in raw mode
			var buf [1]byte
			n, readErr := ui.stdin.Read(buf[:])

			if readErr != nil {
				return 0, readErr
			}

			if n == 0 {
				continue // Try again if no bytes read
			}

			return rune(buf[0]), nil
		}
	}

	// Use buffered reader for non-TTY
	r, _, err := reader.ReadRune()
	return r, err
}

// Extract <...> placeholders from a string
func extractPlaceholders(s string) []string {
	var res []string
	start := -1
	for i, c := range s {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start != -1 {
			res = append(res, s[start:i])
			start = -1
		}
	}
	return res
}

// addCommandToWorkflow stages a command for interactive workflow selection.
func (h *KeyHandler) addCommandToWorkflow(cmdTemplate string) {
	if !h.stageWorkflowCommand(cmdTemplate) {
		return
	}
	h.ui.state.ClearInput()
	h.openWorkflowSelection()
}

func (h *KeyHandler) stageWorkflowCommand(cmdTemplate string) bool {
	parts := strings.Fields(cmdTemplate)
	if len(parts) == 0 {
		return false
	}

	h.ui.pendingWorkflowTemplate = cmdTemplate
	h.ui.pendingWorkflowCommand = parts[0]
	h.ui.pendingWorkflowArgs = append([]string(nil), parts[1:]...)
	return true
}

func (h *KeyHandler) openWorkflowSelection() {
	summaries := h.ui.listWorkflows()
	initialIndex := 0
	activeID := h.ui.workflowMgr.GetActiveID()
	for idx, summary := range summaries {
		if summary.ID == activeID {
			initialIndex = idx
			break
		}
	}
	h.ui.state.ActivateWorkflowSelection(initialIndex)
	h.ui.state.EnterContext(ContextWorkflowSelection)
	h.ui.ensureWorkflowListSelection()
}

func (h *KeyHandler) moveWorkflowSelection(delta int) {
	options := len(h.ui.listWorkflows()) + 1
	if options <= 0 {
		return
	}
	h.ui.state.AdjustWorkflowSelection(delta, options)
}

func (h *KeyHandler) cancelWorkflowSelection() {
	if h.ui.state.workflowSelectionActive {
		h.ui.state.ExitContext()
	}
	h.ui.state.DeactivateWorkflowSelection()
	h.ui.clearPendingWorkflowSelection()
	h.ui.ensureWorkflowListSelection()
}

func (h *KeyHandler) moveWorkflowList(delta int) {
	summaries := h.ui.listWorkflows()
	if len(summaries) == 0 {
		return
	}
	h.ui.state.SetWorkflowListIndex(h.ui.state.workflowListIndex+delta, len(summaries))
}

func (h *KeyHandler) currentWorkflowSummary() (WorkflowSummary, bool) {
	summaries := h.ui.listWorkflows()
	if len(summaries) == 0 {
		return WorkflowSummary{}, false
	}
	idx := h.ui.state.workflowListIndex
	if idx < 0 {
		idx = 0
	}
	if idx >= len(summaries) {
		idx = len(summaries) - 1
	}
	h.ui.state.SetWorkflowListIndex(idx, len(summaries))
	return summaries[idx], true
}

func (h *KeyHandler) createWorkflowFromView() {
	newID := h.ui.workflowMgr.CreateWorkflow()
	h.ui.workflowMgr.SetActive(newID)
	summaries := h.ui.listWorkflows()
	h.ui.state.SetWorkflowListIndex(len(summaries)-1, len(summaries))
	h.ui.ensureWorkflowListSelection()
	h.ui.updateWorkflowPointer()
	h.ui.write("%s➕ Created workflow #%d%s\n",
		h.ui.colors.BrightGreen,
		newID,
		h.ui.colors.Reset)
}

func (h *KeyHandler) copyWorkflowFromView() {
	summary, ok := h.currentWorkflowSummary()
	if !ok {
		h.ui.write("%sNo workflow available to copy%s\n",
			h.ui.colors.BrightYellow,
			h.ui.colors.Reset)
		return
	}

	workflow, exists := h.ui.workflowMgr.GetWorkflow(summary.ID)
	if !exists || workflow == nil {
		h.ui.write("%sUnable to locate workflow #%d%s\n",
			h.ui.colors.BrightRed,
			summary.ID,
			h.ui.colors.Reset)
		return
	}

	templates := workflowToTemplates(workflow)
	if len(templates) == 0 {
		h.ui.write("%sWorkflow #%d has no steps to copy%s\n",
			h.ui.colors.BrightYellow,
			summary.ID,
			h.ui.colors.Reset)
		return
	}

	baseName := strings.TrimSpace(summary.Name)
	if baseName == "" {
		baseName = fmt.Sprintf("Workflow #%d", summary.ID)
	}
	copyName := fmt.Sprintf("%s (copy)", baseName)

	newID, err := h.ui.workflowMgr.CreateWorkflowFromTemplates(copyName, templates)
	if err != nil {
		h.ui.write("%sFailed to copy workflow: %v%s\n",
			h.ui.colors.BrightRed,
			err,
			h.ui.colors.Reset)
		return
	}

	h.ui.workflowMgr.SetActive(newID)
	h.ui.ensureWorkflowListSelection()
	h.ui.updateWorkflowPointer()
	h.ui.write("%s📄 Copied to workflow #%d:%s %s%s%s\n\n",
		h.ui.colors.BrightGreen,
		newID,
		h.ui.colors.Reset,
		h.ui.colors.BrightWhite+h.ui.colors.Bold,
		copyName,
		h.ui.colors.Reset)
}

type workflowSavePlan struct {
	Name       string
	Templates  []string
	Normalized []string
}

func (h *KeyHandler) saveWorkflowFromView() {
	summary, workflow, ok := h.workflowSaveContext()
	if !ok {
		return
	}

	plan, ok := h.prepareWorkflowSave(summary, workflow)
	if !ok {
		return
	}

	if !h.persistWorkflowToConfig(plan) {
		return
	}

	h.registerConfigWorkflowInView(plan)
	h.ui.write("%s💾 Saved workflow to config as %s%s%s\n\n",
		h.ui.colors.BrightGreen,
		h.ui.colors.BrightWhite+h.ui.colors.Bold,
		plan.Name,
		h.ui.colors.Reset)
}

func (h *KeyHandler) workflowSaveContext() (*WorkflowSummary, *Workflow, bool) {
	summary, ok := h.currentWorkflowSummary()
	if !ok {
		h.ui.write("%sNo workflow selected to save%s\n",
			h.ui.colors.BrightYellow,
			h.ui.colors.Reset)
		return nil, nil, false
	}

	if summary.ReadOnly {
		h.ui.write("%sWorkflow #%d is already provided by config. Use copy to duplicate it.%s\n",
			h.ui.colors.BrightYellow,
			summary.ID,
			h.ui.colors.Reset)
		return nil, nil, false
	}

	if h.ui.configMgr == nil || h.ui.config == nil {
		h.ui.write("%sSaving workflows requires write access to configuration.%s\n",
			h.ui.colors.BrightYellow,
			h.ui.colors.Reset)
		return nil, nil, false
	}

	workflow, exists := h.ui.workflowMgr.GetWorkflow(summary.ID)
	if !exists || workflow == nil {
		h.ui.write("%sUnable to locate workflow #%d%s\n",
			h.ui.colors.BrightRed,
			summary.ID,
			h.ui.colors.Reset)
		return nil, nil, false
	}

	return &summary, workflow, true
}

func (h *KeyHandler) prepareWorkflowSave(summary *WorkflowSummary, workflow *Workflow) (*workflowSavePlan, bool) {
	templates := workflowToTemplates(workflow)
	if len(templates) == 0 {
		h.ui.write("%sWorkflow #%d has no steps to save%s\n",
			h.ui.colors.BrightYellow,
			summary.ID,
			h.ui.colors.Reset)
		return nil, false
	}

	name := strings.TrimSpace(summary.Name)
	if name == "" {
		name = fmt.Sprintf("workflow-%d", len(h.ui.config.Workflows)+1)
	}

	normalized := normalizeTemplates(templates)
	for _, existing := range h.ui.config.Workflows {
		if templatesEqual(normalized, normalizeTemplates(existing.Steps)) {
			h.ui.write("%sA workflow with the same steps already exists in config.%s\n",
				h.ui.colors.BrightYellow,
				h.ui.colors.Reset)
			return nil, false
		}
	}

	return &workflowSavePlan{Name: name, Templates: templates, Normalized: normalized}, true
}

func (h *KeyHandler) persistWorkflowToConfig(plan *workflowSavePlan) bool {
	h.ui.config.Workflows = append(h.ui.config.Workflows, config.WorkflowConfig{
		Name:  plan.Name,
		Steps: plan.Templates,
	})

	if err := h.ui.configMgr.Save(); err != nil {
		h.ui.config.Workflows = h.ui.config.Workflows[:len(h.ui.config.Workflows)-1]
		h.ui.write("%sFailed to save workflow to config: %v%s\n",
			h.ui.colors.BrightRed,
			err,
			h.ui.colors.Reset)
		return false
	}

	return true
}

func (h *KeyHandler) registerConfigWorkflowInView(plan *workflowSavePlan) {
	alreadyPresent := false
	for _, existing := range h.ui.listWorkflows() {
		if existing.Source != WorkflowSourceConfig {
			continue
		}
		existingWorkflow, ok := h.ui.workflowMgr.GetWorkflow(existing.ID)
		if !ok || existingWorkflow == nil {
			continue
		}
		if templatesEqual(plan.Normalized, normalizeTemplates(workflowToTemplates(existingWorkflow))) {
			alreadyPresent = true
			break
		}
	}

	if alreadyPresent {
		return
	}

	if _, err := h.ui.workflowMgr.CreateReadOnlyWorkflow(plan.Name, plan.Templates); err != nil {
		h.ui.write("%sSaved to config, but failed to register workflow in view: %v%s\n",
			h.ui.colors.BrightYellow,
			err,
			h.ui.colors.Reset)
		return
	}

	h.ui.ensureWorkflowListSelection()
}

func (h *KeyHandler) finalizeWorkflowSelection() {
	if h.ui.pendingWorkflowCommand == "" {
		h.cancelWorkflowSelection()
		return
	}

	summaries := h.ui.listWorkflows()
	creating := h.ui.state.workflowSelectionIndex >= len(summaries)

	var workflowID int
	if creating {
		workflowID = h.ui.workflowMgr.CreateWorkflow()
	} else {
		idx := h.ui.state.workflowSelectionIndex
		workflowID = summaries[idx].ID
	}

	h.ui.workflowMgr.SetActive(workflowID)
	h.ui.updateWorkflowPointer()
	stepID, err := h.ui.AddToWorkflowByID(
		workflowID,
		h.ui.pendingWorkflowCommand,
		append([]string(nil), h.ui.pendingWorkflowArgs...),
		h.ui.pendingWorkflowTemplate,
	)
	if err != nil {
		h.ui.write("%sFailed to add command to workflow #%d: %v%s\n",
			h.ui.colors.BrightRed,
			workflowID,
			err,
			h.ui.colors.Reset)
		h.cancelWorkflowSelection()
		return
	}

	placeholders := extractPlaceholders(h.ui.pendingWorkflowTemplate)
	h.ui.write("\n%s🎯 Added to workflow #%d!%s\n",
		h.ui.colors.BrightGreen+h.ui.colors.Bold,
		workflowID,
		h.ui.colors.Reset)

	if creating {
		h.ui.write("%s  Created new workflow #%d%s\n",
			h.ui.colors.BrightBlack,
			workflowID,
			h.ui.colors.Reset)
	}

	if len(placeholders) > 0 {
		h.ui.write("%s  Step %d: %s%s%s %s(will prompt for: %v)%s\n",
			h.ui.colors.BrightCyan,
			stepID,
			h.ui.colors.BrightWhite+h.ui.colors.Bold,
			h.ui.pendingWorkflowTemplate,
			h.ui.colors.Reset,
			h.ui.colors.BrightYellow,
			placeholders,
			h.ui.colors.Reset)
	} else {
		h.ui.write("%s  Step %d: %s%s%s\n",
			h.ui.colors.BrightCyan,
			stepID,
			h.ui.colors.BrightWhite+h.ui.colors.Bold,
			h.ui.pendingWorkflowTemplate,
			h.ui.colors.Reset)
	}

	h.ui.write("%s  Press 'Ctrl+t' to manage workflows, or continue adding commands%s\n\n",
		h.ui.colors.BrightBlack,
		h.ui.colors.Reset)

	h.cancelWorkflowSelection()
}

// clearWorkflow clears all steps from workflow
func (h *KeyHandler) clearWorkflow() {
	if h.ui.state.showWorkflow {
		summaries := h.ui.listWorkflows()
		if len(summaries) == 0 {
			h.ui.write("%sNo workflows available to delete.%s\n",
				h.ui.colors.BrightYellow, h.ui.colors.Reset)
			return
		}
		idx := h.ui.state.workflowListIndex
		if idx < 0 {
			idx = 0
		}
		if idx >= len(summaries) {
			idx = len(summaries) - 1
		}
		targetID := summaries[idx].ID
		newActive, ok := h.ui.workflowMgr.DeleteWorkflow(targetID)
		if !ok {
			h.ui.write("%sUnable to delete workflow #%d%s\n",
				h.ui.colors.BrightRed, targetID, h.ui.colors.Reset)
			return
		}
		h.ui.workflowMgr.SetActive(newActive)
		h.ui.ensureWorkflowListSelection()
		h.ui.updateWorkflowPointer()
		h.ui.write("%s🗑️  Deleted workflow #%d%s\n",
			h.ui.colors.BrightYellow,
			targetID,
			h.ui.colors.Reset)
		return
	}

	if err := h.ui.ClearWorkflowByID(h.ui.workflowMgr.GetActiveID()); err != nil {
		h.ui.write("%sFailed to clear workflow: %v%s\n",
			h.ui.colors.BrightRed, err, h.ui.colors.Reset)
		return
	}
	h.ui.write("%s🧹 Active workflow cleared%s\n", h.ui.colors.BrightYellow, h.ui.colors.Reset)
}

// executeWorkflow executes the current workflow
func (h *KeyHandler) executeWorkflow(oldState *term.State) (bool, []string) {
	workflowID := h.executionTargetWorkflowID()
	workflow := h.workflowByID(workflowID)
	if workflow == nil || workflow.IsEmpty() {
		h.ui.write("%sWorkflow is empty. Add some steps first!%s\n",
			h.ui.colors.BrightRed, h.ui.colors.Reset)
		return true, nil
	}

	h.restoreTerminalBeforeWorkflow(oldState)
	clearScreen(h.ui.stdout)

	if err := h.ui.ExecuteWorkflowByID(workflowID); err != nil {
		if errors.Is(err, ErrWorkflowCanceled) {
			h.handleSoftCancel(oldState)
			return true, nil
		}
		fmt.Printf("\n❌ Workflow execution failed: %v\n", err)
	} else {
		fmt.Printf("\n✨ Workflow preserved for reuse. Press 'Ctrl+t' to view or modify.\n")
	}

	return false, []string{"ggc", InteractiveWorkflowCommand}
}

func (h *KeyHandler) executionTargetWorkflowID() int {
	if !h.ui.state.showWorkflow {
		return h.ui.workflowMgr.GetActiveID()
	}
	summaries := h.ui.listWorkflows()
	if len(summaries) == 0 {
		return h.ui.workflowMgr.GetActiveID()
	}
	idx := h.ui.state.workflowListIndex
	if idx < 0 {
		idx = 0
	}
	if idx >= len(summaries) {
		idx = len(summaries) - 1
	}
	return summaries[idx].ID
}

func (h *KeyHandler) workflowByID(workflowID int) *Workflow {
	h.ui.workflowMgr.SetActive(workflowID)
	h.ui.updateWorkflowPointer()
	workflow, exists := h.ui.workflowMgr.GetWorkflow(workflowID)
	if !exists {
		return nil
	}
	return workflow
}

func (h *KeyHandler) restoreTerminalBeforeWorkflow(oldState *term.State) {
	if oldState == nil {
		return
	}
	if f, ok := h.ui.stdin.(*os.File); ok {
		if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
			h.ui.writeError("failed to restore terminal state: %v", err)
		}
	}
}
