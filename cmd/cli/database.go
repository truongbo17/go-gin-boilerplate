package cli

import (
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// WithDatabase runs a command with a database connection managed by the caller.
type WithDatabase func(cmd *cobra.Command, run func(*gorm.DB) error) error
