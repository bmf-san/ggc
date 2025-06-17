package cmd

import (
	"fmt"

	"github.com/bmf-san/ggc/git"
)

type Fetcher struct {
	FetchPrune func() error
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		FetchPrune: git.FetchPrune,
	}
}

func (f *Fetcher) Fetch(args []string) {
	if len(args) > 0 && args[0] == "--prune" {
		err := f.FetchPrune()
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
	ShowFetchHelp()
}

func ShowFetchHelp() {
	fmt.Println("Usage: ggc fetch --prune")
}

// 旧インターフェース維持用ラッパー
// func Fetch(args []string) {
// 	NewFetcher().Fetch(args)
// }
