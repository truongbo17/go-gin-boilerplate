package cli

import (
	"errors"
	"net/mail"

	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"gorm.io/gorm"
)

func init() {
	CreateUserCmd.Flags().StringP("username", "u", "", "username")
	CreateUserCmd.Flags().StringP("password", "p", "", "password")
	CreateUserCmd.Flags().StringP("email", "e", "", "email")
	CreateUserCmd.Flags().Bool("admin", false, "grant the admin role for initial setup")
	_ = CreateUserCmd.MarkFlagRequired("username")
	_ = CreateUserCmd.MarkFlagRequired("password")
	_ = CreateUserCmd.MarkFlagRequired("email")
}

var CreateUserCmd = &cobra.Command{
	Use:   "create_user",
	Short: "Create a user; pass --admin only for trusted operators",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")
		admin, _ := cmd.Flags().GetBool("admin")
		if username == "" || email == "" || len(password) < 8 {
			return errors.New("username, email and password of at least 8 characters are required")
		}
		if _, err := mail.ParseAddress(email); err != nil {
			return err
		}
		hashed, err := (&services.AuthService{}).GeneratePassword(password)
		if err != nil {
			return err
		}
		return database.DB.Transaction(func(tx *gorm.DB) error {
			user := models.User{Username: username, Email: email, Password: string(hashed), Status: enums.StatusActive}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			if !admin {
				return nil
			}
			role := models.Role{Slug: "admin", Name: "Administrator"}
			if err := tx.Where("slug = ?", role.Slug).FirstOrCreate(&role).Error; err != nil {
				return err
			}
			assignment := models.UserRole{UserID: user.ID, RoleID: role.ID}
			return tx.Where("user_id = ? AND role_id = ?", assignment.UserID, assignment.RoleID).FirstOrCreate(&assignment).Error
		})
	},
}
