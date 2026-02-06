package keybindings

// CreateEmacsProfile returns the Emacs-style keybinding profile.
// Based on GNU Emacs standard keybindings with authentic Emacs behavior.
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
