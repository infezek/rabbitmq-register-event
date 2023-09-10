package server

import (
	"github.com/spf13/cobra"
)

func Run() {

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(
		runStart(),
		runFind(),
	)
	rootCmd.Execute()
}
