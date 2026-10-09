package cli

import (
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/internal/migrations"
	"gorm.io/gorm"
)

func NewMigrate(withDatabase WithDatabase) *cobra.Command {
	return &cobra.Command{
		Use:     "migrate",
		Short:   "Run database migrations",
		Example: "ggb migrate",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return withDatabase(cmd, func(db *gorm.DB) error {
				if err := migrations.Run(db); err != nil {
					return err
				}
				cmd.Println("Migration successful")
				return nil
			})
		},
	}
}
