package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Build command that prints hello",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("hello")
    },
}

func init() {
    rootCmd.AddCommand(buildCmd)
}
