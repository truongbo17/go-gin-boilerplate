package cli

import (
	"errors"
	"io"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	authrepository "github.com/truongbo17/go-gin-boilerplate/internal/repository/auth"
	"gorm.io/gorm"
)

func NewCreateUser(withDatabase WithDatabase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_user",
		Short: "Create a user; pass --admin only for trusted operators",
		RunE: func(cmd *cobra.Command, _ []string) error {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			fromStdin, _ := cmd.Flags().GetBool("password-stdin")
			if fromStdin {
				if password != "" {
					return errors.New("use either --password or --password-stdin")
				}
				input, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), 1025))
				if err != nil || len(input) > 1024 {
					return errors.New("cannot read password from standard input")
				}
				password = strings.TrimRight(string(input), "\r\n")
			}
			email, _ := cmd.Flags().GetString("email")
			admin, _ := cmd.Flags().GetBool("admin")
			if username == "" || email == "" {
				return errors.New("username and email are required")
			}
			passwordLength := utf8.RuneCountInString(password)
			if passwordLength < 8 || passwordLength > 64 || len(password) > 72 {
				return errors.New("password must be 8 to 64 characters and at most 72 bytes")
			}
			if len(username) > 50 || len(email) > 100 {
				return errors.New("username must be at most 50 bytes and email at most 100 bytes")
			}
			address, err := mail.ParseAddress(email)
			if err != nil {
				return err
			}
			if address.Address != email {
				return errors.New("email must be a plain address")
			}
			return withDatabase(cmd, func(db *gorm.DB) error {
				userService := services.NewUserService(authrepository.NewUserRepository(db))
				return userService.CreateUser(cmd.Context(), types.ProvisionUserInput{
					Username: username,
					Email:    email,
					Password: password,
					Admin:    admin,
				})
			})
		},
	}
	cmd.Flags().StringP("username", "u", "", "username")
	cmd.Flags().StringP("password", "p", "", "password")
	cmd.Flags().Bool("password-stdin", false, "read password from standard input")
	cmd.Flags().StringP("email", "e", "", "email")
	cmd.Flags().Bool("admin", false, "grant the admin role for initial setup")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")
	return cmd
}
