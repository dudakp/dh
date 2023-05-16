package cmd

import (
	"dh/internal/sqlexecutor"
	"dh/pkg/executor"
	"github.com/spf13/cobra"
	"log"
)

var (
	query       string
	sqlExecutor *sqlexecutor.SqlExecutorService
)

func init() {
	sqlExecutor = sqlexecutor.NewSqlExecutorService(executor.SqlExecutorConfig{})
	qhCommand.
		Flags().
		StringVarP(&query, "query", "q", "", "SQL query")
}

var qhCommand = &cobra.Command{
	Use:   "qh",
	Short: "query helper - collection of SQL queries",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Print(args)
		log.Print(sqlExecutor.ConfigIsEmpty())
	},
}
