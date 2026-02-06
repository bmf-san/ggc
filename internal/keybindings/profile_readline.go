package keybindings

// CreateReadlineProfile returns the GNU Readline-style keybinding profile
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

