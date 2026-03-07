package cmd

import "github.com/spf13/cobra"

// NewRootCmd builds and returns the root cobra command with all subcommands wired.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "worky",
		Short: "CLI for creating and building worky workshops",
	}

	root.AddCommand(
		NewInitCmd(),
		NewNewCmd(),
		NewBuildCmd(),
	)

	return root
}
