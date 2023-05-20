package cmd

import (
	"dh/internal/qh"
	"dh/pkg/executor"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	query       string
	sqlExecutor *qh.SqlExecutorService
)

func init() {
	sqlExecutor = qh.NewSqlExecutorService(executor.SqlExecutorConfig{})
	qhCommand.
		Flags().
		StringVarP(&query, "query", "q", "", "SQL query")
}

var qhCommand = &cobra.Command{
	Use:   "qh",
	Short: "query helper - collection of SQL queries",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if sqlExecutor.ConfigIsEmpty() {
			fmt.Println("qh is not configured, please use c arg to set configuration")
		}
	},
}
