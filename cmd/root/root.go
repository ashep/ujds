package root

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/ashep/go-cfgloader"
	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/config"
	"github.com/ashep/ujds/logger"
	"github.com/ashep/ujds/migration"
	"github.com/ashep/ujds/server"
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

			cfg := config.Config{}
			if err := cfgloader.Load(configPath, &cfg, config.Schema); err != nil {
				l.Fatal().Err(err).Msg("failed to load config")
				return
			}

			if cfg.DB.DSN == "" {
				l.Fatal().Msg("empty db dsn")
				return
			}

			db, err := sql.Open("postgres", cfg.DB.DSN)
			if err != nil {
				l.Fatal().Err(err).Msg("failed to open db")
				return
			}

			if err = db.PingContext(cmd.Context()); err != nil {
				l.Fatal().Err(err).Msg("failed to connect to db")
			}
			l.Debug().Msg("db connection ok")

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

			s := server.New(cfg.Server, api.New(cfg.API, db, l), l.With().Str("pkg", "server").Logger())

			if err := s.Run(cmd.Context()); errors.Is(err, http.ErrServerClosed) {
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
