package cmd

import (
	"testing"
)

func TestNewAddCommitPusher(t *testing.T) {
	pusher := NewAddCommitPusher()
	if pusher == nil {
		t.Fatal("Expected AddCommitPusher, got nil")
	}
	// Basic field checks
	if pusher.execCommand == nil || pusher.inputReader == nil || pusher.outputWriter == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewBrancher(t *testing.T) {
	brancher := NewBrancher()
	if brancher == nil {
		t.Fatal("Expected Brancher, got nil")
	}
	// Basic field checks
	if brancher.gitClient == nil || brancher.execCommand == nil || brancher.inputReader == nil ||
		brancher.outputWriter == nil || brancher.helper == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewCommitter(t *testing.T) {
	committer := NewCommitter()
	if committer == nil {
		t.Fatal("Expected Committer, got nil")
	}
	// Basic field checks
	if committer.gitClient == nil || committer.outputWriter == nil || committer.helper == nil || committer.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewCommitPusher(t *testing.T) {
	pusher := NewCommitPusher()
	if pusher == nil {
		t.Fatal("Expected CommitPusher, got nil")
	}
	// Basic field checks
	if pusher.execCommand == nil || pusher.inputReader == nil || pusher.outputWriter == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewFetcher(t *testing.T) {
	fetcher := NewFetcher()
	if fetcher == nil {
		t.Fatal("Expected Fetcher, got nil")
	}
	// Basic field checks
	if fetcher.outputWriter == nil || fetcher.helper == nil || fetcher.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewDiffer(t *testing.T) {
	differ := NewDiffer()
	if differ == nil {
		t.Fatal("Expected Differ, got nil")
	}
	// Basic field checks
	if differ.outputWriter == nil || differ.helper == nil || differ.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("Expected Logger, got nil")
	}
	// Basic field checks
	if logger.gitClient == nil || logger.outputWriter == nil || logger.helper == nil || logger.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewPuller(t *testing.T) {
	puller := NewPuller()
	if puller == nil {
		t.Fatal("Expected Puller, got nil")
	}
	// Basic field checks
	if puller.gitClient == nil || puller.outputWriter == nil || puller.helper == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewPullRebasePusher(t *testing.T) {
	pusher := NewPullRebasePusher()
	if pusher == nil {
		t.Fatal("Expected PullRebasePusher, got nil")
	}
	// Basic field checks
	if pusher.gitClient == nil || pusher.outputWriter == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewPusher(t *testing.T) {
	pusher := NewPusher()
	if pusher == nil {
		t.Fatal("Expected Pusher, got nil")
	}
	// Basic field checks
	if pusher.gitClient == nil || pusher.outputWriter == nil || pusher.helper == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewRebaser(t *testing.T) {
	rebaser := NewRebaser()
	if rebaser == nil {
		t.Fatal("Expected Rebaser, got nil")
	}
	// Basic field checks
	if rebaser.outputWriter == nil || rebaser.helper == nil || rebaser.execCommand == nil || rebaser.inputReader == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewRemoteer(t *testing.T) {
	remoteer := NewRemoteer()
	if remoteer == nil {
		t.Fatal("Expected Remoteer, got nil")
	}
	// Basic field checks
	if remoteer.execCommand == nil || remoteer.outputWriter == nil || remoteer.helper == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewResetter(t *testing.T) {
	resetter := NewResetter()
	if resetter == nil {
		t.Fatal("Expected Resetter, got nil")
	}
	// Basic field checks
	if resetter.outputWriter == nil || resetter.helper == nil || resetter.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewResetCleaner(t *testing.T) {
	cleaner := NewResetCleaner()
	if cleaner == nil {
		t.Fatal("Expected ResetCleaner, got nil")
	}
	// Basic field checks
	if cleaner.outputWriter == nil || cleaner.helper == nil || cleaner.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewStasher(t *testing.T) {
	stasher := NewStasher()
	if stasher == nil {
		t.Fatal("Expected Stasher, got nil")
	}
	// Basic field checks
	if stasher.outputWriter == nil || stasher.helper == nil || stasher.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewConfigureer(t *testing.T) {
	configureer := NewConfigureer()
	if configureer == nil {
		t.Fatal("Expected configureer, got nil")
	}
	// Basic field checks
	if configureer.outputWriter == nil || configureer.helper == nil || configureer.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewStatuseer(t *testing.T) {
	statuseer := NewStatuseer()
	if statuseer == nil {
		t.Fatal("Expected Statuseer, got nil")
	}
	// Basic field checks
	if statuseer.outputWriter == nil || statuseer.helper == nil || statuseer.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewTagger(t *testing.T) {
	tagger := NewTagger()
	if tagger == nil {
		t.Fatal("Expected Tagger, got nil")
	}
	// Basic field checks
	if tagger.outputWriter == nil || tagger.helper == nil || tagger.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewVersioneer(t *testing.T) {
	versioneer := NewVersioneer()
	if versioneer == nil {
		t.Fatal("Expected Versioneer, got nil")
	}
	// Basic field checks
	if versioneer.outputWriter == nil || versioneer.helper == nil || versioneer.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewStashPullPopper(t *testing.T) {
	popper := NewStashPullPopper()
	if popper == nil {
		t.Fatal("Expected StashPullPopper, got nil")
	}
	// Basic field checks
	if popper.outputWriter == nil || popper.helper == nil || popper.execCommand == nil {
		t.Error("Expected all fields to be initialized")
	}
}

func TestNewCmd_Constructor(t *testing.T) {
	cmd := NewCmd()
	if cmd == nil {
		t.Fatal("Expected Cmd, got nil")
	}
	// Basic field checks - just verify the main components exist
	if cmd.gitClient == nil || cmd.outputWriter == nil || cmd.helper == nil {
		t.Error("Expected core fields to be initialized")
	}
	// Verify all command handlers are initialized
	if cmd.brancher == nil || cmd.committer == nil || cmd.logger == nil ||
		cmd.puller == nil || cmd.pusher == nil || cmd.resetter == nil ||
		cmd.cleaner == nil || cmd.pullRebasePusher == nil || cmd.adder == nil ||
		cmd.remoteer == nil || cmd.rebaser == nil || cmd.stasher == nil ||
		cmd.commitPusher == nil || cmd.addCommitPusher == nil || cmd.completer == nil ||
		cmd.fetcher == nil || cmd.stashPullPopper == nil || cmd.resetCleaner == nil ||
		cmd.statuseer == nil || cmd.differ == nil || cmd.tagger == nil || cmd.versioneer == nil ||
		cmd.configureer == nil {
		t.Error("Expected all command handlers to be initialized")
	}
}
