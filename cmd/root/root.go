package root

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ashep/datapimp/app"
	"github.com/ashep/datapimp/config"
	"github.com/ashep/datapimp/logger"
	"github.com/ashep/datapimp/migration"
	"github.com/spf13/cobra"
)

var (
	debugMode  bool
	configPath string
	migUp      bool
	migDown    bool
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			l := logger.New(debugMode)

			cfg, err := config.ParseFromPath(configPath)
			if err != nil {
				l.Fatal().Err(err).Msg("failed to load config")
				return
			}

			if cfg.DB.DSN == "" {
				cfg.DB.DSN = "postgres://postgres:postgres@localhost/datapimp?sslmode=disable"
			}

			db, err := sql.Open("postgres", cfg.DB.DSN)
			if err != nil {
				l.Fatal().Err(err).Msg("failed to connect to db")
				return
			}

			if migUp {
				if err := migration.Up(db); err != nil {
					l.Fatal().Err(err).Msg("failed to apply migrations")
				}
				return
			}

			if migDown {
				if err := migration.Down(db); err != nil {
					l.Fatal().Err(err).Msg("failed to revert migrations")
				}
				return
			}

			if err := app.Run(cmd.Context(), cfg, db, l); errors.Is(err, http.ErrServerClosed) {
				l.Info().Msg("server stopped")
			} else if err != nil {
				l.Error().Err(err).Msg("")
			}
		},
	}

	cmd.Flags().BoolVar(&migUp, "migrate-up", false, "apply database migrations")
	cmd.Flags().BoolVar(&migDown, "migrate-down", false, "revert database migrations")

	cmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug mode")
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "path to the config file")

	return cmd
}
