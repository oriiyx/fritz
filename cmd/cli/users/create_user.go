package users

import (
	"errors"
	"net/mail"

	"github.com/oriiyx/fritz/cmd/cli/config"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/spf13/cobra"
)

func NewCreateUserCmd(deps *config.Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "create <email> <password>",
		Short: "Create a new user",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			email := args[0]
			password := args[1]

			if _, err := mail.ParseAddress(email); err != nil {
				return errors.New("email not correct format")
			}

			if password == "" {
				return errors.New("missing password")
			}

			// Now you have access to deps.Queries, deps.DB, deps.Logger
			_, err := deps.Queries.CreateUser(cmd.Context(), db.CreateUserParams{
				Email:    email,
				Password: password,
			})
			if err != nil {
				deps.Logger.Error().Err(err).Msg("Failed to create user")
				return err
			}

			deps.Logger.Info().Str("email", email).Msg("User created successfully")
			return nil
		},
	}
}
