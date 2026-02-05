package keybindings

import "fmt"

// Profile represents built-in keybinding profiles that users can select
type Profile string

// Supported keybinding profiles.
const (
	ProfileDefault  Profile = "default"  // Current default behavior (backward compatible)
	ProfileEmacs    Profile = "emacs"    // Emacs-style bindings (Ctrl-based, modeless)
	ProfileVi       Profile = "vi"       // Vi-style bindings (modal concepts adapted for CLI)
	ProfileReadline Profile = "readline" // GNU Readline standard bindings
)

// String returns the string representation of a Profile
func (p Profile) String() string {
	return string(p)
}

// IsValid checks if a Profile value is valid
func (p Profile) IsValid() bool {
	switch p {
	case ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline:
		return true
	default:
		return false
	}
}

// Context represents different UI states that can have specific keybindings
type Context string

// Available contexts for interactive UI states.
const (
	ContextGlobal  Context = "global"  // Always active (reserved keys like Ctrl+C)
	ContextInput   Context = "input"   // When typing/editing the search query
	ContextResults Context = "results" // When navigating through filtered results
	ContextSearch  Context = "search"  // When fuzzy search is active (combines input + results)
)

// String returns the string representation of a Context
func (c Context) String() string {
	return string(c)
}

// IsValid checks if a Context value is valid
func (c Context) IsValid() bool {
	switch c {
	case ContextGlobal, ContextInput, ContextResults, ContextSearch:
		return true
	default:
		return false
	}
}

// GetAllProfiles returns a list of all valid profiles
func GetAllProfiles() []Profile {
	return []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}
}

// GetAllContexts returns a list of all valid contexts
func GetAllContexts() []Context {
	return []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
}

// KeyBindingProfile defines keybindings for a complete profile
type KeyBindingProfile struct {
	Name        string                             // Profile name (e.g., "emacs")
	Description string                             // Human-readable description
	Global      map[string][]KeyStroke             // Global keybindings (always active)
	Contexts    map[Context]map[string][]KeyStroke // Context-specific keybindings
}

// NewKeyBindingProfile creates a new profile with initialized maps
func NewKeyBindingProfile(name, description string) *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        name,
		Description: description,
		Global:      make(map[string][]KeyStroke),
		Contexts:    make(map[Context]map[string][]KeyStroke),
	}
}

// SetGlobalBinding sets a global keybinding (active in all contexts)
func (kbp *KeyBindingProfile) SetGlobalBinding(action string, keystrokes []KeyStroke) {
	if kbp.Global == nil {
		kbp.Global = make(map[string][]KeyStroke)
	}
	kbp.Global[action] = keystrokes
}

// SetContextBinding sets a context-specific keybinding
func (kbp *KeyBindingProfile) SetContextBinding(context Context, action string, keystrokes []KeyStroke) {
	if kbp.Contexts == nil {
		kbp.Contexts = make(map[Context]map[string][]KeyStroke)
	}
	if kbp.Contexts[context] == nil {
		kbp.Contexts[context] = make(map[string][]KeyStroke)
	}
	kbp.Contexts[context][action] = keystrokes
}

// GetBinding returns the keybinding for an action in a specific context
// Falls back to global bindings if not found in context
func (kbp *KeyBindingProfile) GetBinding(context Context, action string) ([]KeyStroke, bool) {
	// Try context-specific first
	if contextMap, exists := kbp.Contexts[context]; exists {
		if keystrokes, exists := contextMap[action]; exists {
			return keystrokes, true
		}
	}

	// Fall back to global
	if keystrokes, exists := kbp.Global[action]; exists {
		return keystrokes, true
	}

	return nil, false
}

// GetAllActions returns all action names defined in this profile
func (kbp *KeyBindingProfile) GetAllActions() []string {
	actionSet := make(map[string]bool)

	// Add global actions
	for action := range kbp.Global {
		actionSet[action] = true
	}

	// Add context-specific actions
	for _, contextMap := range kbp.Contexts {
		for action := range contextMap {
			actionSet[action] = true
		}
	}

	// Convert to slice
	actions := make([]string, 0, len(actionSet))
	for action := range actionSet {
		actions = append(actions, action)
	}

	return actions
}

// Clone creates a deep copy of the profile
func (kbp *KeyBindingProfile) Clone() *KeyBindingProfile {
	clone := NewKeyBindingProfile(kbp.Name, kbp.Description)

	// Clone global bindings
	for action, keystrokes := range kbp.Global {
		clonedKeystrokes := make([]KeyStroke, len(keystrokes))
		copy(clonedKeystrokes, keystrokes)
		clone.Global[action] = clonedKeystrokes
	}

	// Clone context bindings
	for context, contextMap := range kbp.Contexts {
		clone.Contexts[context] = make(map[string][]KeyStroke)
		for action, keystrokes := range contextMap {
			clonedKeystrokes := make([]KeyStroke, len(keystrokes))
			copy(clonedKeystrokes, keystrokes)
			clone.Contexts[context][action] = clonedKeystrokes
		}
	}

	return clone
}

// CreateDefaultProfile returns the default keybinding profile (legacy compatible)
func CreateDefaultProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Default",
		Description: "Default keybindings compatible with legacy behavior",
		Global:      make(map[string][]KeyStroke),
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"soft_cancel": {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				"delete_word":       {NewCtrlKeyStroke('w')},
				"clear_line":        {NewCtrlKeyStroke('u')},
				"delete_to_end":     {NewCtrlKeyStroke('k')},
				"move_to_beginning": {NewCtrlKeyStroke('a')},
				"move_to_end":       {NewCtrlKeyStroke('e')},
			},
			ContextResults: {
				"move_up":              {NewCtrlKeyStroke('p')},
				"move_down":            {NewCtrlKeyStroke('n')},
				"add_to_workflow":      {NewTabKeyStroke()},
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},
				"clear_workflow":       {NewCharKeyStroke('c')},
			},
			ContextSearch: {
				"move_up":              {NewCtrlKeyStroke('p')},
				"move_down":            {NewCtrlKeyStroke('n')},
				"add_to_workflow":      {NewTabKeyStroke()},
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},
				"clear_workflow":       {NewCharKeyStroke('c')},
			},
		},
	}
}

// CreateEmacsProfile returns the Emacs-style keybinding profile
// Based on GNU Emacs standard keybindings with authentic Emacs behavior
func CreateEmacsProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Emacs",
		Description: "Comprehensive Emacs-style keybindings with authentic GNU Emacs behavior",
		Global: map[string][]KeyStroke{
			// Core Emacs global bindings
			"quit":                {NewCtrlKeyStroke('g')},                        // C-g keyboard-quit
			"help":                {NewCtrlKeyStroke('h')},                        // C-h help-command
			"universal_argument":  {NewCtrlKeyStroke('u')},                        // C-u universal-argument
			"exchange_point_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x (chord)
			"suspend":             {NewCtrlKeyStroke('z')},                        // C-z suspend-frame
		},
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"quit":               {NewCtrlKeyStroke('g')},
				"help":               {NewCtrlKeyStroke('h')},
				"universal_argument": {NewCtrlKeyStroke('u')},
				"suspend":            {NewCtrlKeyStroke('z')},
				"soft_cancel":        {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				// Character-level movement
				"forward_char":  {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char": {NewCtrlKeyStroke('b')}, // C-b backward-char
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line

				// Word-level movement
				"forward_word":  {NewAltKeyStroke('f', "")}, // M-f forward-word
				"backward_word": {NewAltKeyStroke('b', "")}, // M-b backward-word

				// Line-level movement
				"beginning_of_line": {NewCtrlKeyStroke('a')}, // C-a beginning-of-line
				"end_of_line":       {NewCtrlKeyStroke('e')}, // C-e end-of-line
				"move_to_beginning": {NewCtrlKeyStroke('a')}, // Alias for compatibility
				"move_to_end":       {NewCtrlKeyStroke('e')}, // Alias for compatibility

				// Deletion and killing
				"delete_char":          {NewCtrlKeyStroke('d')},                        // C-d delete-char
				"backward_delete_char": {NewCtrlKeyStroke('h')},                        // C-h backward-delete-char
				"kill_line":            {NewCtrlKeyStroke('k')},                        // C-k kill-line
				"kill_word":            {NewAltKeyStroke('d', "")},                     // M-d kill-word
				"backward_kill_word":   {NewAltKeyStroke(127, "backspace")},            // M-DEL backward-kill-word
				"unix_line_discard":    {NewCtrlKeyStroke('u')},                        // C-u unix-line-discard
				"kill_whole_line":      {NewCtrlKeyStroke('s'), NewCtrlKeyStroke('k')}, // C-S-k kill-whole-line
				"delete_word":          {NewAltKeyStroke('d', "")},                     // Alias for kill-word
				"clear_line":           {NewCtrlKeyStroke('u')},                        // Alias for unix-line-discard
				"delete_to_end":        {NewCtrlKeyStroke('k')},                        // Alias for kill-line

				// Search and replace
				"isearch_forward":  {NewCtrlKeyStroke('s')},    // C-s isearch-forward
				"isearch_backward": {NewCtrlKeyStroke('r')},    // C-r isearch-backward
				"query_replace":    {NewAltKeyStroke('%', "")}, // M-% query-replace

				// Case operations
				"upcase_word":     {NewAltKeyStroke('u', "")}, // M-u upcase-word
				"downcase_word":   {NewAltKeyStroke('l', "")}, // M-l downcase-word
				"capitalize_word": {NewAltKeyStroke('c', "")}, // M-c capitalize-word
				"transpose_chars": {NewCtrlKeyStroke('t')},    // C-t transpose-chars
				"transpose_words": {NewAltKeyStroke('t', "")}, // M-t transpose-words

				// Yank and kill ring
				"yank":                {NewCtrlKeyStroke('y')},    // C-y yank
				"yank_pop":            {NewAltKeyStroke('y', "")}, // M-y yank-pop
				"copy_region_as_kill": {NewAltKeyStroke('w', "")}, // M-w copy-region-as-kill
				"kill_region":         {NewCtrlKeyStroke('w')},    // C-w kill-region

				// Mark and region
				"set_mark_command":    {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark-command
				"exchange_point_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x exchange-point-mark

				// Buffer and file operations (adapted for CLI)
				"save_buffer":      {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('s')}, // C-x C-s save-buffer
				"find_file":        {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('f')}, // C-x C-f find-file
				"switch_to_buffer": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('b')}, // C-x C-b switch-to-buffer

				// Miscellaneous
				"quoted_insert":           {NewCtrlKeyStroke('q')},     // C-q quoted-insert
				"recenter_top_bottom":     {NewCtrlKeyStroke('l')},     // C-l recenter-top-bottom
				"just_one_space":          {NewAltKeyStroke(' ', "")},  // M-SPC just-one-space
				"delete_horizontal_space": {NewAltKeyStroke('\\', "")}, // M-\ delete-horizontal-space
			},
			ContextResults: {
				// Navigation in results (Emacs-style list navigation)
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"move_up":       {NewCtrlKeyStroke('p')}, // Alias
				"move_down":     {NewCtrlKeyStroke('n')}, // Alias
				"backward_char": {NewCtrlKeyStroke('b')}, // C-b backward-char
				"forward_char":  {NewCtrlKeyStroke('f')}, // C-f forward-char

				// Scrolling (Emacs page movement)
				"scroll_up":           {NewAltKeyStroke('v', "")}, // M-v scroll-up
				"scroll_down":         {NewCtrlKeyStroke('v')},    // C-v scroll-down
				"beginning_of_buffer": {NewAltKeyStroke('<', "")}, // M-< beginning-of-buffer
				"end_of_buffer":       {NewAltKeyStroke('>', "")}, // M-> end-of-buffer

				// Selection and marking
				"set_mark_command":  {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark-command
				"mark_whole_buffer": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('h')}, // C-x h mark-whole-buffer

				// Search in results
				"isearch_forward":  {NewCtrlKeyStroke('s')}, // C-s isearch-forward
				"isearch_backward": {NewCtrlKeyStroke('r')}, // C-r isearch-backward

				// Execute/select
				"execute": {NewCtrlKeyStroke('m')}, // C-m (Enter equivalent)
				"select":  {NewCtrlKeyStroke('m')}, // Alias

				// Workflow operations (adapted for Emacs style)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})}, // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},      // C-t
				"clear_workflow":       {NewAltKeyStroke('c', "")},   // M-c clear
			},
			ContextSearch: {
				// Search-specific Emacs bindings
				"isearch_forward":         {NewCtrlKeyStroke('s')},       // C-s isearch-forward
				"isearch_backward":        {NewCtrlKeyStroke('r')},       // C-r isearch-backward
				"isearch_repeat_forward":  {NewCtrlKeyStroke('s')},       // C-s (repeat)
				"isearch_repeat_backward": {NewCtrlKeyStroke('r')},       // C-r (repeat)
				"isearch_yank_word":       {NewCtrlKeyStroke('w')},       // C-w isearch-yank-word
				"isearch_yank_line":       {NewCtrlKeyStroke('y')},       // C-y isearch-yank-line
				"isearch_delete_char":     {NewCtrlKeyStroke('h')},       // C-h isearch-delete-char
				"isearch_abort":           {NewCtrlKeyStroke('g')},       // C-g isearch-abort
				"isearch_exit":            {NewRawKeyStroke([]byte{13})}, // RET isearch-exit

				// Navigation while searching
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line
				"move_up":       {NewCtrlKeyStroke('p')}, // Alias
				"move_down":     {NewCtrlKeyStroke('n')}, // Alias

				// Case sensitivity toggle
				"isearch_toggle_case_fold": {NewAltKeyStroke('c', "")}, // M-c toggle case sensitivity
				"isearch_toggle_regexp":    {NewAltKeyStroke('r', "")}, // M-r toggle regexp mode

				// Workflow operations (search context)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})}, // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},      // C-t
				"clear_workflow":       {NewAltKeyStroke('x', "")},   // M-x clear (avoiding conflict with M-c)
			},
		},
	}
}

// CreateViProfile returns the Vi-style keybinding profile (adapted for CLI context)
// Implements Vi modal editing concepts adapted for command-line interface
func CreateViProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Vi",
		Description: "Vi-style modal keybindings adapted for command-line interface with insert and normal modes",
		Global: map[string][]KeyStroke{
			// Core Vi global bindings
			"quit":          {NewCtrlKeyStroke('c')},             // Keep standard quit (like :q!)
			"command_mode":  {NewRawKeyStroke([]byte{27})},       // ESC - enter command mode
			"force_quit":    {NewRawKeyStroke([]byte{'Z', 'Q'})}, // ZQ - quit without saving
			"save_and_quit": {NewRawKeyStroke([]byte{'Z', 'Z'})}, // ZZ - save and quit
		},
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"quit":          {NewCtrlKeyStroke('c')},
				"command_mode":  {NewRawKeyStroke([]byte{27})},
				"force_quit":    {NewRawKeyStroke([]byte{'Z', 'Q'})},
				"save_and_quit": {NewRawKeyStroke([]byte{'Z', 'Z'})},
				"soft_cancel":   {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				// Vi INSERT MODE bindings (when editing input)
				// In Vi, insert mode is similar to normal editor behavior

				// Basic movement (limited in insert mode)
				"move_to_beginning": {NewCtrlKeyStroke('a')}, // C-a move to beginning
				"move_to_end":       {NewCtrlKeyStroke('e')}, // C-e move to end
				"forward_char":      {NewCtrlKeyStroke('l')}, // C-l move right
				"backward_char":     {NewCtrlKeyStroke('h')}, // C-h move left (also backspace)

				// Deletion (insert mode)
				"delete_word":          {NewCtrlKeyStroke('w')}, // C-w delete word backward
				"delete_line":          {NewCtrlKeyStroke('u')}, // C-u delete line
				"clear_line":           {NewCtrlKeyStroke('u')}, // Alias
				"delete_to_end":        {NewCtrlKeyStroke('k')}, // C-k delete to end of line
				"backward_delete_char": {NewCtrlKeyStroke('h')}, // C-h backspace

				// Insert mode specific
				"insert_at_beginning": {NewRawKeyStroke([]byte{'I'})}, // I - insert at line beginning
				"insert_at_end":       {NewRawKeyStroke([]byte{'A'})}, // A - insert at line end
				"open_line_below":     {NewRawKeyStroke([]byte{'o'})}, // o - open new line below
				"open_line_above":     {NewRawKeyStroke([]byte{'O'})}, // O - open new line above

				// Exit insert mode
				"escape_to_normal": {NewRawKeyStroke([]byte{27})}, // ESC - to normal mode

				// Vi-style completion and registers
				"complete_word":  {NewCtrlKeyStroke('n')}, // C-n word completion
				"complete_prev":  {NewCtrlKeyStroke('p')}, // C-p previous completion
				"literal_insert": {NewCtrlKeyStroke('v')}, // C-v literal character insert
			},
			ContextResults: {
				// Vi NORMAL MODE bindings (when navigating results)
				// This is where Vi really shines with single-key navigation

				// Basic movement (hjkl)
				"move_left":  {NewRawKeyStroke([]byte{'h'})}, // h - move left
				"move_down":  {NewRawKeyStroke([]byte{'j'})}, // j - move down
				"move_up":    {NewRawKeyStroke([]byte{'k'})}, // k - move up
				"move_right": {NewRawKeyStroke([]byte{'l'})}, // l - move right

				// Alternative movement for compatibility
				"move_down_alt": {NewCtrlKeyStroke('n')}, // C-n alternative
				"move_up_alt":   {NewCtrlKeyStroke('p')}, // C-p alternative

				// Word movement
				"forward_word":      {NewRawKeyStroke([]byte{'w'})}, // w - next word
				"backward_word":     {NewRawKeyStroke([]byte{'b'})}, // b - previous word
				"end_word":          {NewRawKeyStroke([]byte{'e'})}, // e - end of word
				"forward_word_big":  {NewRawKeyStroke([]byte{'W'})}, // W - next WORD
				"backward_word_big": {NewRawKeyStroke([]byte{'B'})}, // B - previous WORD
				"end_word_big":      {NewRawKeyStroke([]byte{'E'})}, // E - end of WORD

				// Line movement
				"first_char":        {NewRawKeyStroke([]byte{'^'})}, // ^ - first non-blank character
				"beginning_of_line": {NewRawKeyStroke([]byte{'0'})}, // 0 - beginning of line
				"end_of_line":       {NewRawKeyStroke([]byte{'$'})}, // $ - end of line

				// Screen movement
				"top_of_screen":    {NewRawKeyStroke([]byte{'H'})}, // H - top of screen
				"middle_of_screen": {NewRawKeyStroke([]byte{'M'})}, // M - middle of screen
				"bottom_of_screen": {NewRawKeyStroke([]byte{'L'})}, // L - bottom of screen

				// Buffer movement
				"first_line": {NewRawKeyStroke([]byte{'g', 'g'})}, // gg - first line
				"last_line":  {NewRawKeyStroke([]byte{'G'})},      // G - last line
				"goto_line":  {NewRawKeyStroke([]byte{':'})},      // : - command mode (go to line)

				// Scrolling
				"scroll_down":      {NewCtrlKeyStroke('f')}, // C-f - page down
				"scroll_up":        {NewCtrlKeyStroke('b')}, // C-b - page up
				"scroll_down_half": {NewCtrlKeyStroke('d')}, // C-d - half page down
				"scroll_up_half":   {NewCtrlKeyStroke('u')}, // C-u - half page up
				"scroll_line_down": {NewCtrlKeyStroke('e')}, // C-e - scroll down one line
				"scroll_line_up":   {NewCtrlKeyStroke('y')}, // C-y - scroll up one line

				// Search and navigation
				"search_forward":       {NewRawKeyStroke([]byte{'/'})}, // / - search forward
				"search_backward":      {NewRawKeyStroke([]byte{'?'})}, // ? - search backward
				"search_next":          {NewRawKeyStroke([]byte{'n'})}, // n - next search match
				"search_previous":      {NewRawKeyStroke([]byte{'N'})}, // N - previous search match
				"search_word_forward":  {NewRawKeyStroke([]byte{'*'})}, // * - search word under cursor forward
				"search_word_backward": {NewRawKeyStroke([]byte{'#'})}, // # - search word under cursor backward

				// Marks and jumps
				"set_mark":       {NewRawKeyStroke([]byte{'m'})},  // m{a-z} - set mark
				"goto_mark":      {NewRawKeyStroke([]byte{'\''})}, // '{a-z} - goto mark
				"goto_mark_line": {NewRawKeyStroke([]byte{'`'})},  // `{a-z} - goto mark exact position
				"jump_back":      {NewCtrlKeyStroke('o')},         // C-o - jump back
				"jump_forward":   {NewCtrlKeyStroke('i')},         // C-i - jump forward

				// Selection and execution
				"select":           {NewRawKeyStroke([]byte{13})},  // Enter - select current item
				"execute":          {NewRawKeyStroke([]byte{13})},  // Alias
				"visual_mode":      {NewRawKeyStroke([]byte{'v'})}, // v - visual mode
				"visual_line_mode": {NewRawKeyStroke([]byte{'V'})}, // V - visual line mode

				// Repeat and undo (adapted for CLI)
				"repeat_last": {NewRawKeyStroke([]byte{'.'})}, // . - repeat last action
				"undo":        {NewRawKeyStroke([]byte{'u'})}, // u - undo
				"redo":        {NewCtrlKeyStroke('r')},        // C-r - redo

				// Enter insert mode from results
				"insert_mode":         {NewRawKeyStroke([]byte{'i'})}, // i - insert mode
				"insert_after":        {NewRawKeyStroke([]byte{'a'})}, // a - insert after cursor
				"insert_at_end":       {NewRawKeyStroke([]byte{'A'})}, // A - insert at line end
				"insert_at_beginning": {NewRawKeyStroke([]byte{'I'})}, // I - insert at line beginning

				// Workflow operations (Vi normal mode style)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},   // Tab
				"toggle_workflow_view": {NewRawKeyStroke([]byte{'W'})}, // W - workflow view (capital W)
				"clear_workflow":       {NewRawKeyStroke([]byte{'D'})}, // D - delete/clear workflow
			},
			ContextSearch: {
				// Vi search mode bindings (when in / or ? search)
				// Similar to insert mode but with search-specific commands

				// Basic navigation
				"move_up":       {NewRawKeyStroke([]byte{'k'})}, // k - move up in results
				"move_down":     {NewRawKeyStroke([]byte{'j'})}, // j - move down in results
				"move_up_alt":   {NewCtrlKeyStroke('p')},        // C-p alternative
				"move_down_alt": {NewCtrlKeyStroke('n')},        // C-n alternative

				// Search navigation
				"search_next":     {NewRawKeyStroke([]byte{'n'})}, // n - next match
				"search_previous": {NewRawKeyStroke([]byte{'N'})}, // N - previous match
				"search_repeat":   {NewRawKeyStroke([]byte{13})},  // Enter - accept search
				"search_abort":    {NewRawKeyStroke([]byte{27})},  // ESC - abort search

				// Edit search term
				"delete_word":  {NewCtrlKeyStroke('w')}, // C-w delete word
				"clear_search": {NewCtrlKeyStroke('u')}, // C-u clear search line
				"delete_char":  {NewCtrlKeyStroke('h')}, // C-h delete character

				// Search modes
				"case_sensitive_toggle": {NewRawKeyStroke([]byte{'\\', 'c'})}, // \c - toggle case sensitivity
				"regex_mode_toggle":     {NewRawKeyStroke([]byte{'\\', 'v'})}, // \v - very magic mode
				"literal_mode_toggle":   {NewRawKeyStroke([]byte{'\\', 'V'})}, // \V - very nomagic mode

				// History (search command history)
				"search_history_up":   {NewCtrlKeyStroke('p')}, // C-p - previous search
				"search_history_down": {NewCtrlKeyStroke('n')}, // C-n - next search

				// Workflow operations (Vi search mode)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},   // Tab
				"toggle_workflow_view": {NewRawKeyStroke([]byte{'W'})}, // W - workflow view
				"clear_workflow":       {NewRawKeyStroke([]byte{'D'})}, // D - delete/clear workflow
			},
		},
	}
}

// CreateReadlineProfile returns the GNU Readline compatible keybinding profile
// Based on GNU Readline library defaults providing bash-like experience
func CreateReadlineProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Readline",
		Description: "Comprehensive GNU Readline compatible keybindings for authentic bash-like CLI experience",
		Global: map[string][]KeyStroke{
			// Core Readline global bindings
			"abort":        {NewCtrlKeyStroke('g')}, // C-g abort
			"bell":         {NewCtrlKeyStroke('g')}, // C-g bell (same as abort)
			"clear_screen": {NewCtrlKeyStroke('l')}, // C-l clear-screen
		},
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"abort":        {NewCtrlKeyStroke('g')},
				"clear_screen": {NewCtrlKeyStroke('l')},
				"soft_cancel":  {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				// Character Movement (GNU Readline standard)
				"forward_char":      {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char":     {NewCtrlKeyStroke('b')}, // C-b backward-char
				"move_to_beginning": {NewCtrlKeyStroke('a')}, // C-a beginning-of-line
				"move_to_end":       {NewCtrlKeyStroke('e')}, // C-e end-of-line
				"beginning_of_line": {NewCtrlKeyStroke('a')}, // Alias
				"end_of_line":       {NewCtrlKeyStroke('e')}, // Alias

				// Word Movement
				"forward_word":  {NewAltKeyStroke('f', "")}, // M-f forward-word
				"backward_word": {NewAltKeyStroke('b', "")}, // M-b backward-word

				// Line Navigation
				"next_line":        {NewCtrlKeyStroke('n')}, // C-n next-history
				"previous_line":    {NewCtrlKeyStroke('p')}, // C-p previous-history
				"previous_history": {NewCtrlKeyStroke('p')}, // Alias
				"next_history":     {NewCtrlKeyStroke('n')}, // Alias

				// Character Deletion
				"delete_char":          {NewCtrlKeyStroke('d')}, // C-d delete-char
				"backward_delete_char": {NewCtrlKeyStroke('h')}, // C-h backward-delete-char (backspace)

				// Word Deletion
				"kill_word":          {NewAltKeyStroke('d', "")},          // M-d kill-word
				"backward_kill_word": {NewAltKeyStroke(127, "backspace")}, // M-DEL backward-kill-word
				"unix_word_rubout":   {NewCtrlKeyStroke('w')},             // C-w unix-word-rubout
				"delete_word":        {NewCtrlKeyStroke('w')},             // Alias for compatibility

				// Line Killing and Yanking
				"kill_line":         {NewCtrlKeyStroke('k')},                        // C-k kill-line
				"unix_line_discard": {NewCtrlKeyStroke('u')},                        // C-u unix-line-discard
				"kill_whole_line":   {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('k')}, // C-x C-k kill-whole-line
				"clear_line":        {NewCtrlKeyStroke('u')},                        // Alias
				"delete_to_end":     {NewCtrlKeyStroke('k')},                        // Alias

				// Yank and Kill Ring
				"yank":          {NewCtrlKeyStroke('y')},    // C-y yank
				"yank_pop":      {NewAltKeyStroke('y', "")}, // M-y yank-pop
				"yank_nth_arg":  {NewAltKeyStroke('.', "")}, // M-. yank-nth-arg (yank last arg)
				"yank_last_arg": {NewAltKeyStroke('_', "")}, // M-_ yank-last-arg

				// Transposition
				"transpose_chars": {NewCtrlKeyStroke('t')},    // C-t transpose-chars
				"transpose_words": {NewAltKeyStroke('t', "")}, // M-t transpose-words

				// Case Manipulation
				"upcase_word":     {NewAltKeyStroke('u', "")}, // M-u upcase-word
				"downcase_word":   {NewAltKeyStroke('l', "")}, // M-l downcase-word
				"capitalize_word": {NewAltKeyStroke('c', "")}, // M-c capitalize-word

				// History Operations
				"reverse_search_history":  {NewCtrlKeyStroke('r')},    // C-r reverse-search-history
				"forward_search_history":  {NewCtrlKeyStroke('s')},    // C-s forward-search-history
				"history_search_backward": {NewAltKeyStroke('p', "")}, // M-p history-search-backward
				"history_search_forward":  {NewAltKeyStroke('n', "")}, // M-n history-search-forward
				"beginning_of_history":    {NewAltKeyStroke('<', "")}, // M-< beginning-of-history
				"end_of_history":          {NewAltKeyStroke('>', "")}, // M-> end-of-history

				// Completion
				"complete":             {NewRawKeyStroke([]byte{9})}, // TAB complete
				"possible_completions": {NewAltKeyStroke('?', "")},   // M-? possible-completions
				"insert_completions":   {NewAltKeyStroke('*', "")},   // M-* insert-completions
				"complete_filename":    {NewAltKeyStroke('/', "")},   // M-/ complete-filename
				"complete_username":    {NewAltKeyStroke('~', "")},   // M-~ complete-username
				"complete_variable":    {NewAltKeyStroke('$', "")},   // M-$ complete-variable
				"complete_hostname":    {NewAltKeyStroke('@', "")},   // M-@ complete-hostname

				// Numeric Arguments
				"digit_argument":     {NewAltKeyStroke('0', "")}, // M-0 through M-9 digit-argument
				"universal_argument": {NewCtrlKeyStroke('u')},    // C-u universal-argument

				// Miscellaneous
				"quoted_insert":           {NewCtrlKeyStroke('v')},                        // C-v quoted-insert
				"tab_insert":              {NewAltKeyStroke('\t', "")},                    // M-TAB tab-insert
				"tilde_expand":            {NewAltKeyStroke('&', "")},                     // M-& tilde-expand
				"set_mark":                {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark
				"exchange_point_and_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x exchange-point-and-mark

				// Editing Commands
				"overwrite_mode": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('o')}, // C-x C-o overwrite-mode
				"undo":           {NewCtrlKeyStroke('_')},                        // C-_ undo
				"revert_line":    {NewAltKeyStroke('r', "")},                     // M-r revert-line

				// Shell Integration
				"glob_complete_word":   {NewAltKeyStroke('g', "")},                     // M-g glob-complete-word
				"glob_expand_word":     {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('*')}, // C-x * glob-expand-word
				"glob_list_expansions": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('g')}, // C-x g glob-list-expansions

				// Line Editing
				"accept_line": {NewRawKeyStroke([]byte{13})}, // RET accept-line
				"newline":     {NewRawKeyStroke([]byte{10})}, // LFD newline

				// Special Characters
				"self_insert":           {NewRawKeyStroke([]byte{' '})},                     // printable chars self-insert
				"bracketed_paste_begin": {NewRawKeyStroke([]byte{27, 91, 50, 48, 48, 126})}, // bracketed paste mode

				// Macro Operations
				"start_kbd_macro":     {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('(')}, // C-x ( start-kbd-macro
				"end_kbd_macro":       {NewCtrlKeyStroke('x'), NewCtrlKeyStroke(')')}, // C-x ) end-kbd-macro
				"call_last_kbd_macro": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('e')}, // C-x e call-last-kbd-macro

				// Advanced Readline Features
				"dump_functions": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('f')}, // C-x C-f dump-functions
				"dump_variables": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('v')}, // C-x C-v dump-variables
				"dump_macros":    {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('m')}, // C-x C-m dump-macros

				// Menu Complete (bash 4.0+)
				"menu_complete":          {NewAltKeyStroke('\t', "")}, // M-TAB menu-complete
				"menu_complete_backward": {NewAltKeyStroke('\\', "")}, // M-\ menu-complete-backward

				// Delete and Space Manipulation
				"delete_horizontal_space": {NewAltKeyStroke('\\', "")}, // M-\ delete-horizontal-space
				"just_one_space":          {NewAltKeyStroke(' ', "")},  // M-SPC just-one-space
			},
			ContextResults: {
				// Navigation in results using Readline conventions
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"move_up":       {NewCtrlKeyStroke('p')}, // Alias
				"move_down":     {NewCtrlKeyStroke('n')}, // Alias

				// Horizontal movement
				"forward_char":  {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char": {NewCtrlKeyStroke('b')}, // C-b backward-char

				// Page movement
				"scroll_up":   {NewAltKeyStroke('v', "")}, // M-v scroll-up
				"scroll_down": {NewCtrlKeyStroke('v')},    // C-v scroll-down

				// List navigation
				"beginning_of_buffer": {NewAltKeyStroke('<', "")}, // M-< beginning-of-buffer
				"end_of_buffer":       {NewAltKeyStroke('>', "")}, // M-> end-of-buffer

				// Selection
				"accept_line": {NewRawKeyStroke([]byte{13})}, // RET accept-line
				"select":      {NewRawKeyStroke([]byte{13})}, // Alias

				// Search in results
				"reverse_search_history": {NewCtrlKeyStroke('r')}, // C-r reverse-search
				"forward_search_history": {NewCtrlKeyStroke('s')}, // C-s forward-search

				// Mark and selection
				"set_mark":                {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark
				"exchange_point_and_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x exchange-point-and-mark

				// Workflow operations (Readline style)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},                   // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('w')}, // C-x C-w workflow
				"clear_workflow":       {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('c')}, // C-x C-c clear
			},
			ContextSearch: {
				// Search mode using Readline search conventions
				"search_forward":  {NewCtrlKeyStroke('s')},       // C-s search-forward
				"search_backward": {NewCtrlKeyStroke('r')},       // C-r search-backward
				"search_abort":    {NewCtrlKeyStroke('g')},       // C-g abort-search
				"search_accept":   {NewRawKeyStroke([]byte{13})}, // RET accept-search

				// Navigation in search
				"move_up":   {NewCtrlKeyStroke('p')}, // C-p previous-match
				"move_down": {NewCtrlKeyStroke('n')}, // C-n next-match

				// Edit search string
				"delete_char":          {NewCtrlKeyStroke('d')}, // C-d delete-char
				"backward_delete_char": {NewCtrlKeyStroke('h')}, // C-h backward-delete-char
				"kill_line":            {NewCtrlKeyStroke('k')}, // C-k kill-line
				"unix_line_discard":    {NewCtrlKeyStroke('u')}, // C-u unix-line-discard
				"delete_word":          {NewCtrlKeyStroke('w')}, // C-w delete-word

				// Search string movement
				"forward_char":      {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char":     {NewCtrlKeyStroke('b')}, // C-b backward-char
				"beginning_of_line": {NewCtrlKeyStroke('a')}, // C-a beginning-of-line
				"end_of_line":       {NewCtrlKeyStroke('e')}, // C-e end-of-line

				// Search history
				"search_history_up":   {NewCtrlKeyStroke('p')}, // C-p previous-search
				"search_history_down": {NewCtrlKeyStroke('n')}, // C-n next-search

				// Search completion
				"complete":             {NewRawKeyStroke([]byte{9})}, // TAB complete-search
				"possible_completions": {NewAltKeyStroke('?', "")},   // M-? possible-completions

				// Yank into search
				"yank":          {NewCtrlKeyStroke('y')},    // C-y yank
				"yank_last_arg": {NewAltKeyStroke('.', "")}, // M-. yank-last-arg

				// Workflow operations (search context)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},                   // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('w')}, // C-x C-w workflow
				"clear_workflow":       {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('c')}, // C-x C-c clear
			},
		},
	}
}

// RegisterBuiltinProfiles registers all built-in profiles with the resolver
func RegisterBuiltinProfiles(resolver *KeyBindingResolver) {
	resolver.RegisterProfile(ProfileDefault, CreateDefaultProfile())
	resolver.RegisterProfile(ProfileEmacs, CreateEmacsProfile())
	resolver.RegisterProfile(ProfileVi, CreateViProfile())
	resolver.RegisterProfile(ProfileReadline, CreateReadlineProfile())
}

// GetAllProfilesBuiltin returns all available profile names
func GetAllProfilesBuiltin() []Profile {
	return []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}
}

// GetProfileDescription returns a description for a profile
func GetProfileDescription(profile Profile) string {
	switch profile {
	case ProfileDefault:
		return "Default keybindings compatible with legacy behavior"
	case ProfileEmacs:
		return "Comprehensive Emacs-style keybindings with authentic GNU Emacs behavior"
	case ProfileVi:
		return "Vi-style modal keybindings adapted for command-line interface with insert and normal modes"
	case ProfileReadline:
		return "Comprehensive GNU Readline compatible keybindings for authentic bash-like CLI experience"
	default:
		return "Unknown profile"
	}
}

// ValidateProfile validates a keybinding profile for consistency and completeness
func ValidateProfile(profile *KeyBindingProfile) error { //nolint:revive // performs exhaustive validation checks
	if profile == nil {
		return fmt.Errorf("profile is nil")
	}

	if profile.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if profile.Description == "" {
		return fmt.Errorf("profile description cannot be empty")
	}

	if profile.Contexts == nil {
		return fmt.Errorf("profile contexts cannot be nil")
	}

	// Validate that profile has required contexts
	requiredContexts := []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
	for _, requiredCtx := range requiredContexts {
		if _, exists := profile.Contexts[requiredCtx]; !exists {
			return fmt.Errorf("profile missing required context: %s", requiredCtx)
		}
	}

	// Validate that each context has at least basic navigation bindings
	if inputBindings, exists := profile.Contexts[ContextInput]; exists {
		requiredInputActions := []string{"move_to_beginning", "move_to_end", "delete_word", "clear_line"}
		for _, action := range requiredInputActions {
			if _, hasAction := inputBindings[action]; !hasAction {
				return fmt.Errorf("profile input context missing required action: %s", action)
			}
		}
	}

	if resultsBindings, exists := profile.Contexts[ContextResults]; exists {
		requiredResultsActions := []string{"move_up", "move_down"}
		for _, action := range requiredResultsActions {
			if _, hasAction := resultsBindings[action]; !hasAction {
				return fmt.Errorf("profile results context missing required action: %s", action)
			}
		}
	}

	// Validate KeyStroke consistency
	for contextName, contextBindings := range profile.Contexts {
		for action, keystrokes := range contextBindings {
			if len(keystrokes) == 0 {
				return fmt.Errorf("profile %s context %s action %s has no keystrokes", profile.Name, contextName, action)
			}
			for i, ks := range keystrokes {
				if err := validateKeyStroke(ks); err != nil {
					return fmt.Errorf("profile %s context %s action %s keystroke %d invalid: %w", profile.Name, contextName, action, i, err)
				}
			}
		}
	}

	return nil
}

// ValidateAllBuiltinProfiles validates all built-in profiles
func ValidateAllBuiltinProfiles() error {
	profiles := map[Profile]func() *KeyBindingProfile{
		ProfileDefault:  CreateDefaultProfile,
		ProfileEmacs:    CreateEmacsProfile,
		ProfileVi:       CreateViProfile,
		ProfileReadline: CreateReadlineProfile,
	}

	for profileName, creator := range profiles {
		profile := creator()
		if err := ValidateProfile(profile); err != nil {
			return fmt.Errorf("built-in profile %s validation failed: %w", profileName, err)
		}
	}

	return nil
}

// GetProfileStatistics returns statistics about a profile's keybinding coverage
func GetProfileStatistics(profile *KeyBindingProfile) map[string]interface{} {
	stats := make(map[string]interface{})

	if profile == nil {
		return stats
	}

	// Count total bindings
	totalBindings := 0
	contextStats := make(map[Context]int)

	for context, bindings := range profile.Contexts {
		count := len(bindings)
		contextStats[context] = count
		totalBindings += count
	}

	// Count global bindings
	globalBindings := 0
	if profile.Global != nil {
		globalBindings = len(profile.Global)
	}

	stats["profile_name"] = profile.Name
	stats["description"] = profile.Description
	stats["total_context_bindings"] = totalBindings
	stats["global_bindings"] = globalBindings
	stats["context_breakdown"] = contextStats
	stats["contexts_defined"] = len(profile.Contexts)

	// Calculate keystroke type distribution
	keystrokeTypes := make(map[KeyStrokeKind]int)
	for _, bindings := range profile.Contexts {
		for _, keystrokes := range bindings {
			for _, ks := range keystrokes {
				keystrokeTypes[ks.Kind]++
			}
		}
	}
	stats["keystroke_types"] = keystrokeTypes

	return stats
}

// CompareProfiles compares two profiles and returns differences
func CompareProfiles(profile1, profile2 *KeyBindingProfile) map[string]interface{} { //nolint:revive // comparison builds rich analysis report
	comparison := make(map[string]interface{})

	if profile1 == nil || profile2 == nil {
		comparison["error"] = "one or both profiles are nil"
		return comparison
	}

	comparison["profile1_name"] = profile1.Name
	comparison["profile2_name"] = profile2.Name

	// Compare contexts
	contexts1 := make(map[Context]bool)
	contexts2 := make(map[Context]bool)

	for ctx := range profile1.Contexts {
		contexts1[ctx] = true
	}
	for ctx := range profile2.Contexts {
		contexts2[ctx] = true
	}

	var uniqueToProfile1, uniqueToProfile2, sharedContexts []Context
	for ctx := range contexts1 {
		if contexts2[ctx] {
			sharedContexts = append(sharedContexts, ctx)
		} else {
			uniqueToProfile1 = append(uniqueToProfile1, ctx)
		}
	}
	for ctx := range contexts2 {
		if !contexts1[ctx] {
			uniqueToProfile2 = append(uniqueToProfile2, ctx)
		}
	}

	comparison["unique_to_profile1"] = uniqueToProfile1
	comparison["unique_to_profile2"] = uniqueToProfile2
	comparison["shared_contexts"] = sharedContexts

	// Compare action coverage in shared contexts
	actionComparison := make(map[Context]map[string]interface{})
	for _, ctx := range sharedContexts {
		bindings1 := profile1.Contexts[ctx]
		bindings2 := profile2.Contexts[ctx]

		actions1 := make(map[string]bool)
		actions2 := make(map[string]bool)

		for action := range bindings1 {
			actions1[action] = true
		}
		for action := range bindings2 {
			actions2[action] = true
		}

		var uniqueActions1, uniqueActions2, sharedActions []string
		for action := range actions1 {
			if actions2[action] {
				sharedActions = append(sharedActions, action)
			} else {
				uniqueActions1 = append(uniqueActions1, action)
			}
		}
		for action := range actions2 {
			if !actions1[action] {
				uniqueActions2 = append(uniqueActions2, action)
			}
		}

		actionComparison[ctx] = map[string]interface{}{
			"unique_to_profile1": uniqueActions1,
			"unique_to_profile2": uniqueActions2,
			"shared_actions":     sharedActions,
		}
	}

	comparison["action_comparison"] = actionComparison
	return comparison
}
