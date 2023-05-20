package cmd

import (
	"dh/internal/qh"
	"dh/pkg/executor"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"os"
)

var (
	setConfig       bool
	executorService *qh.SqlExecutorService
)

func init() {
	executorService = qh.NewSqlExecutorService(executor.SqlExecutorConfig{})
	qhCommand.
		Flags().
		BoolVarP(&setConfig, "conf", "c", false, "set config")
}

var qhCommand = &cobra.Command{
	Use:   "qh",
	Short: "query helper - collection of SQL queries",
	Args:  cobra.MinimumNArgs(1),
	Run:   runQh,
}

func runQh(cmd *cobra.Command, args []string) {
	if executorService.ConfigIsEmpty() && !setConfig {
		fmt.Println("qh is not configured, please use c arg to set configuration")
	}
	if setConfig {
		if _, err := tea.NewProgram(qh.NewViewModel(executorService)).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
		return
	}
	executorService.Run(args[0])
}
