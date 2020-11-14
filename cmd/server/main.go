package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of go-service-generator
	Version = "0"

	// Branch is the branch this binary was built from
	Branch = "0"

	// Commit is the commit this binary was built from
	Commit = "0"

	// BuildTime is the time this binary was built
	BuildTime = ""
)

// rootCmd is the root command for the app.
var rootCmd = &cobra.Command{
	Use:   "go-template",
	Short: "go-template service",
}

// Execute rootCmd.
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
