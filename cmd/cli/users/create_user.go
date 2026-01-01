package users

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/manifoldco/promptui"
	"github.com/oriiyx/fritz/cmd/cli/config"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/spf13/cobra"
)

func NewCreateUserCmd(deps *config.Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "create [email] [password]",
		Short: "Create a new Fritz user account",
		Long: `Create a new user account in the Fritz PIM system.

The email must be a valid email address and will be used for authentication.
The password will be securely hashed before storage.

If email or password are not provided, you will be prompted for them.`,
		Example: `  # Create interactively (will prompt for email and password)
  fritz users create

  # Create with email and password
  fritz users create admin@example.com securePassword123

  # Create with a strong password
  fritz users create developer@company.com 'MyS3cur3P@ssw0rd!'`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var email, password string
			var err error

			// Get email (from args or prompt)
			if len(args) >= 1 {
				email = args[0]
			} else {
				emailPrompt := promptui.Prompt{
					Label: "Email",
					Validate: func(input string) error {
						if _, err := mail.ParseAddress(input); err != nil {
							return errors.New("invalid email format")
						}
						return nil
					},
				}
				email, err = emailPrompt.Run()
				if err != nil {
					return fmt.Errorf("email prompt cancelled: %w", err)
				}
			}

			// Validate email format
			if _, err := mail.ParseAddress(email); err != nil {
				return fmt.Errorf("invalid email format: %s", email)
			}

			// Get password (from args or prompt)
			if len(args) >= 2 {
				password = args[1]
			} else {
				passwordPrompt := promptui.Prompt{
					Label: "Password",
					Mask:  '*',
					Validate: func(input string) error {
						if len(input) < 8 {
							return errors.New("password must be at least 8 characters")
						}
						return nil
					},
				}
				password, err = passwordPrompt.Run()
				if err != nil {
					return fmt.Errorf("password prompt cancelled: %w", err)
				}
			}

			// Validate password
			if password == "" {
				return errors.New("password cannot be empty")
			}

			if len(password) < 8 {
				deps.Logger.Warn().Msg("Password is shorter than recommended 8 characters")
			}

			// Create the user
			user, err := deps.Queries.CreateUser(cmd.Context(), db.CreateUserParams{
				Email:    email,
				Password: password,
			})
			if err != nil {
				deps.Logger.Error().Err(err).Str("email", email).Msg("Failed to create user")
				return fmt.Errorf("failed to create user: %w", err)
			}

			// Success message
			deps.Logger.Info().
				Str("email", email).
				Str("user_id", user.ID.String()).
				Msg("User created successfully")
			deps.Logger.Info().Msg("\n\n")
			deps.Logger.Info().Msg("✅ User created successfully")
			deps.Logger.Info().Msgf("  Email    : %s", email)
			deps.Logger.Info().Msgf("  User ID  : %s", user.ID.String())

			return nil
		},
	}
}
