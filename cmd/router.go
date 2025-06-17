package cmd

func Complete(args []string) {
	completer := NewCompleter()
	completer.Complete(args)
}
