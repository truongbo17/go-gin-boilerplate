package cli

import "github.com/spf13/cobra"

func NewVersion() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Get the version of Go Gin Base",
		Example: "ggb version",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println("version 0.0.1")
		},
	}
}
