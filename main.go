package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "wetalk",
		Short: "wetalk is a chat application for coders",
	}
	rootCmd.AddCommand(cmdServer, cmdClient)
	rootCmd.Execute()
}
