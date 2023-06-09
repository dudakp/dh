package cmd

import (
	mrh2 "github.com/dudakp/dh/internal/mrh"
	"github.com/dudakp/dh/pkg/executor"
	"github.com/spf13/cobra"
)

var (
	mrhService *mrh2.Mrh
)

func init() {
	mrhService = &mrh2.Mrh{}

	mrhCommand.
		Flags().
		StringVarP(&mrhService.BranchType, "branchType", "t", "feature", "Branch type")
	mrhCommand.
		Flags().
		BoolVarP(&mrhService.Done, "done", "d", false, "Merge request done")
}

var mrhCommand = &cobra.Command{
	Use:   "mrh",
	Short: "merge request helper - tool for simpler local git repository management while doing reviewing merge requests",
	Args:  cobra.MinimumNArgs(1),
	Run:   mrh,
}

func mrh(cmd *cobra.Command, args []string) {
	gitExecutor, err := executor.NewGitExecutor()
	if err != nil {
		logger.Fatalf("unable to create file executor: %s, ", err.Error())
	}
	mrhService.GitExecutor = gitExecutor

	issue := args[0]
	mrhService.Run(issue)
}
