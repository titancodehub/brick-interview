package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Entry Point of Transaction Service",
	Long:  "Entry Point of Transaction Service"}

func Execute() {
	rootCmd.AddCommand(StartRestServer)
	rootCmd.AddCommand(StartWorker)
	err := rootCmd.Execute()
	if err != nil {
		panic("Failed to run the command")
	}
}
