package cmd

import "io"

// Executer is an interface for executing commands.
type Executer interface {
	Add(args []string)
	AddCommitPush()
	Branch(args []string)
	Clean(args []string)
	CleanInteractive()
	Commit(args []string)
	CommitPushInteractive()
	Complete(args []string)
	Fetch(args []string)
	Log(args []string)
	Pull(args []string)
	PullRebasePush()
	Push(args []string)
	Rebase(args []string)
	Remote(args []string)
	Reset(args []string)
	Stash(args []string)
	StashPullPop()
	ShowHelp()
}

// Cmd is a struct that holds all the commands.
type Cmd struct {
	outputWriter     io.Writer
	adder            *Adder
	addCommitPusher  *AddCommitPusher
	brancher         *Brancher
	cleaner          *Cleaner
	committer        *Committer
	commitPusher     *CommitPusher
	completer        *Completer
	fetcher          *Fetcher
	helper           *Helper
	logger           *Logger
	puller           *Puller
	pullRebasePusher *PullRebasePusher
	pusher           *Pusher
	rebaser          *Rebaser
	remoteer         *Remoteer
	resetter         *Resetter
	stasher          *Stasher
}

// NewCmd creates a new Cmd.
func NewCmd(w io.Writer) *Cmd {
	return &Cmd{
		outputWriter:     w,
		adder:            NewAdder(),
		addCommitPusher:  NewAddCommitPusher(),
		brancher:         NewBrancher(),
		cleaner:          NewCleaner(),
		committer:        NewCommitter(),
		commitPusher:     NewCommitPusher(),
		completer:        NewCompleter(),
		fetcher:          NewFetcher(),
		helper:           NewHelper(),
		logger:           NewLogger(),
		puller:           NewPuller(),
		pullRebasePusher: NewPullRebasePusher(),
		pusher:           NewPusher(),
		rebaser:          NewRebaser(),
		remoteer:         NewRemoteer(),
		resetter:         NewResetter(),
		stasher:          NewStasher(),
	}
}

func (c *Cmd) Add(args []string) {
	c.adder.Add(args)
}

func (c *Cmd) AddCommitPush() {
	c.addCommitPusher.AddCommitPush()
}

func (c *Cmd) Branch(args []string) {
	c.brancher.Branch(args)
}

func (c *Cmd) Clean(args []string) {
	c.cleaner.Clean(args)
}

func (c *Cmd) CleanInteractive() {
	c.cleaner.CleanInteractive()
}

func (c *Cmd) Commit(args []string) {
	c.committer.Commit(args)
}

func (c *Cmd) CommitPushInteractive() {
	c.commitPusher.CommitPushInteractive()
}

func (c *Cmd) Complete(args []string) {
	c.completer.Complete(args)
}

func (c *Cmd) Fetch(args []string) {
	c.fetcher.Fetch(args)
}

func (c *Cmd) Log(args []string) {
	c.logger.Log(args)
}

func (c *Cmd) Pull(args []string) {
	c.puller.Pull(args)
}

func (c *Cmd) PullRebasePush() {
	c.pullRebasePusher.PullRebasePush()
}

func (c *Cmd) Push(args []string) {
	c.pusher.Push(args)
}

func (c *Cmd) Rebase(args []string) {
	c.rebaser.Rebase(args)
}

func (c *Cmd) Remote(args []string) {
	c.remoteer.Remote(args)
}

func (c *Cmd) Reset(args []string) {
	c.resetter.Reset(args)
}

func (c *Cmd) Stash(args []string) {
	c.stasher.Stash(args)
}

func (c *Cmd) StashPullPop() {
	// not implemented
}

func (c *Cmd) ShowHelp() {
	c.helper.ShowHelp()
}
