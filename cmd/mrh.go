package cmd

import (
	"dh/pkg/git"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var (
	branchType  string
	gitExecutor *git.Executor
)

func init() {
	rootCmd.AddCommand(mrhCommand)
	mrhCommand.
		Flags().
		StringVarP(&branchType, "branchType", "b", "feature", "Branch type")
	gitExecutor = git.NewGitExecutor(findGitExecutable())
}

var mrhCommand = &cobra.Command{
	Use:   "mrh",
	Short: "merge request helper - tool for simpler local git repository management while doing reviewing merge requests",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gitExecutor.Status()
	},
}

func findGitExecutable() string {
	path, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Unable to find git executable!")
	}
	return path
}
