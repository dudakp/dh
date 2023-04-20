package cmd

import (
	"dh/pkg/config"
	"dh/pkg/git"
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var (
	branchType  string
	done        bool
	gitExecutor *git.Executor
)

func init() {
	mrhCommand.
		Flags().
		StringVarP(&branchType, "branchType", "b", "feature", "Branch type")
	mrhCommand.
		Flags().
		BoolVarP(&done, "done", "d", false, "Merge request done")
	gitExecutor = git.NewGitExecutor(findGitExecutable())
}

var mrhCommand = &cobra.Command{
	Use:   "mrh",
	Short: "merge request helper - tool for simpler local git repository management while doing reviewing merge requests",
	Args:  cobra.MinimumNArgs(1),
	Run:   mrh,
}

func mrh(cmd *cobra.Command, args []string) {
	var err error
	issue := args[0]

	// save changes in stash
	if !done {
		err = errors.Join(gitExecutor.Stash(false))
		err = errors.Join(gitExecutor.Checkout(branchType + "/" + issue))
	} else {
		// TODO: add prompt check (checklist like: did test ran succesfully, will sonar check pass? ...)
		// code review done, checkout back to develop and pop stashed changes
		err = errors.Join(gitExecutor.Checkout("develop"))
		err = errors.Join(gitExecutor.Stash(true))
	}

	if err != nil {
		rollback(err)
	} else {
		if !done {
			config.InfoLog.Print("repository is ready for code review")
		} else {
			config.InfoLog.Print("repository has rolled back to state before code review")
		}
	}
}

func rollback(err error) {
	config.ErrLog.Printf("calling rollback action caused by error: %s", err.Error())
	err = gitExecutor.Stash(true)
	if err != nil {
		config.ErrLog.Fatalf("Error during git stash pop: %s. Please resolver this error manually", err.Error())
	}
}

func findGitExecutable() string {
	path, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Unable to find git executable!")
	}
	return path
}
