package cmd

import "github.com/spf13/cobra"

var (
	query string
)

func init() {
	rootCmd.AddCommand(qhCommand)
	qhCommand.
		Flags().
		StringVarP(&query, "query", "q", "", "SQL query")
}

var qhCommand = &cobra.Command{
	Use:   "qh",
	Short: "query helper - collection of SQL queries",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		print("qh")
	},
}
