package cmd

import "github.com/spf13/cobra"

// NewNewCmd returns the `worky new` parent command group.
func NewNewCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "new",
		Short: "Scaffold new workshop components",
	}

	c.AddCommand(NewChapterCmd())

	return c
}
