package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "itc",
		Short: "running in-toto verifications in a container",
	}

	logLevel string
)

func init() {
	rootCmd.AddCommand(
		newVerifyCmd(),
	)

	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
