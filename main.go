package main

import (
	"github.com/spf13/cobra"
)

type Message struct {
	Name, Msg string
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "wetalk",
		Version: "0.1.0",
		Long: `wetalk is a chatroom application for coders

Example:
  start server: wetalk server -p 8086
  start client：wetalk client ws://127.0.0.1:8086 -u somebody
`,
	}
	rootCmd.AddCommand(cmdServer, cmdClient)
	rootCmd.Execute()
}
