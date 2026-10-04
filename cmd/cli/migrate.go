package cli

import (
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/internal/migrations"
)

var MigrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Run database migrations",
	Example: "ggb migrate",
	Run: func(cmd *cobra.Command, args []string) {
		migrateCmd(args)
	},
}

func migrateCmd(args []string) {
	migrations.Migrate()
}
