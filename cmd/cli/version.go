package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get the version of Go Gin Base",
	Example: "ggb version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version 0.0.1")
	},
}
