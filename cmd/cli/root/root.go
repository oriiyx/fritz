package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/oriiyx/fritz/app/core/utils/env"
	logger2 "github.com/oriiyx/fritz/app/core/utils/logger"
	"github.com/oriiyx/fritz/cmd/cli/config"
	"github.com/oriiyx/fritz/cmd/cli/definitions"
	"github.com/oriiyx/fritz/cmd/cli/users"
	"github.com/oriiyx/fritz/cmd/cli/version"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	jsonOutput  bool
	shortOutput bool
)

const (
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(initConfig)

	deps := initDependencies()
	defer deps.DB.Close()

	// Create root command with dependencies
	rootCmd := newRootCmd(deps)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s✗ %s%s\n", ColorRed, err.Error(), ColorReset)
		os.Exit(1)
	}
}

func newRootCmd(deps *config.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fritz",
		Short: "Fritz - Dynamic Data Framework CLI",
		Long: `Fritz CLI - Command-line interface for managing Fritz system.

Provides tools for user management, database operations, and system administration.
Use 'fritz [command] --help' for more information about a command.`,
		Version: version.GetVersion(),
	}

	// Global flags
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fritz/config.yaml)")

	// Silence errors and usage on error (we handle them ourselves)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	// Improve command suggestions
	cmd.DisableSuggestions = false
	cmd.SuggestionsMinimumDistance = 2

	// Disable default completion command (we can add our own later if needed)
	cmd.CompletionOptions.DisableDefaultCmd = true

	// Hide the help command (help is available via --help flag)
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	// Set version template
	cmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)

	// Add subcommands with injected dependencies
	cmd.AddCommand(users.NewUsersCmd(deps))
	cmd.AddCommand(definitions.NewDefinitionsCmd(deps))
	cmd.AddCommand(newVersionCmd())

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Create .fritz directory if it doesn't exist
		fritzDir := filepath.Join(home, config.Directory)
		if err := os.MkdirAll(fritzDir, 0750); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating %s directory: %v\n", config.Directory, err)
			os.Exit(1)
		}

		// Search config in .fritz directory
		viper.AddConfigPath(fritzDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
}

func initDependencies() *config.Dependencies {
	// Load config
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
	conf := env.New()

	// Initialize logger
	logger, err := logger2.New(true, "var/logs/cli.log")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	dbString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		conf.DB.Host, conf.DB.Username, conf.DB.Password, conf.DB.DBName, conf.DB.Port)

	pool, err := pgxpool.New(context.Background(), dbString)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create database pool")
	}

	queries := db.New(pool)

	return &config.Dependencies{
		DB:      pool,
		Queries: queries,
		Logger:  logger,
	}
}
