package keybindings

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
