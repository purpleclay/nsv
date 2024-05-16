package main

import (
	"fmt"
	"os"

	"github.com/purpleclay/nsv/cmd"
)

var (
	// The current built version
	version = ""
	// The git branch associated with the current built version
	gitBranch = ""
	// The git SHA1 of the commit
	gitCommit = ""
	// The date associated with the current built version
	buildDate = ""
)

func main() {
	err := cmd.Execute(os.Stdout, cmd.BuildDetails{
		Version:   version,
		GitBranch: gitBranch,
		GitCommit: gitCommit,
		Date:      buildDate,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
