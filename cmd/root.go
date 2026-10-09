package cmd

import (
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/cmd/cli"
)

var rootCmd = &cobra.Command{
	Use:   "ggb",
	Short: "Reusable Gin API base",
}

func init() {
	rootCmd.AddCommand(StartServerCmd, StartWorkerCmd, cli.NewVersion(), cli.NewMigrate(withDatabase), cli.NewCreateUser(withDatabase))
}

func Execute() error {
	return rootCmd.Execute()
}
