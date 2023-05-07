package cmd

import (
	"dh/internal"
	"dh/pkg/executor"
	"github.com/spf13/cobra"
)

var mrhService *internal.Mrh

func init() {
	mrhService = &internal.Mrh{
		GitExecutor: executor.NewGitExecutor(),
	}

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
	issue := args[0]
	mrhService.Run(issue)
}
