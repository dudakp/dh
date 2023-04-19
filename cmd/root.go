package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "dh",
	Short: "dh (developer helper) is collection of tools for everyday use",
	Run: func(cmd *cobra.Command, args []string) {
		print("Use command")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
