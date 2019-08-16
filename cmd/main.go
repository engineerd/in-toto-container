package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "itc",
	Short: "running in-toto verifications in a container",
}

func init() {
	rootCmd.AddCommand(
		newVerifyCmd(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
