package cmd

import (
	"dh/internal/qh"
	"dh/pkg/executor"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"os"
)

var (
	setConfig       bool
	executorService *qh.SqlExecutorService
	listQueries     bool
)

func init() {
	qhCommand.
		Flags().
		BoolVarP(&setConfig, "conf", "c", false, "set config")
	qhCommand.
		Flags().
		BoolVarP(&listQueries, "listQueries", "l", false, "List all available queries")
}

var qhCommand = &cobra.Command{
	Use:   "qh",
	Short: "query helper - collection of SQL queries",
	Run:   runQh,
}

func runQh(cmd *cobra.Command, args []string) {
	es, err := qh.NewSqlExecutorService()
	executorService = es
	if err != nil {
		panic(err)
	}
	if executorService.ConfigIsEmpty() && !setConfig {
		fmt.Println("qh is not configured, please use c arg to set configuration")
		return
	}
	if setConfig {
		configViewModel := qh.NewViewModel(executorService)
		if _, err := tea.NewProgram(configViewModel).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
		return
	}
	if listQueries {
		queries := executorService.ListAvailableQueries()
		if len(queries) > 0 {
			fmt.Printf("%d available queries:\n", len(queries))
		} else {
			fmt.Println("no available queries, please configure query dir path")
			return
		}
		for _, query := range queries {
			fmt.Println(fmt.Sprintf("%s", query))
		}
		return
	}
	if len(args) < 1 {
		fmt.Println("missing name of query to be executed as arg")
		return
	}
	var queryData executor.QueryData
	if len(args) > 1 {
		queryData.Arg = args[1]
	}
	queryName := args[0]
	res, err := executorService.Run(queryName, queryData)
	if err != nil {
		if errors.Is(err, executor.QueryNotFound) {
			fmt.Printf("query: %s not found", queryName)
			os.Exit(0)
		}
	}
	resultSetModel := qh.NewResultModel(res)
	if _, err := tea.NewProgram(resultSetModel).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
