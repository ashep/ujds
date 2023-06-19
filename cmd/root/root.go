package root

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/ashep/go-cfgloader"
	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/config"
	"github.com/ashep/ujds/logger"
	"github.com/ashep/ujds/migration"
	"github.com/ashep/ujds/server"
)

var (
	cfgPath string
	migUp   bool
	migDown bool
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			appName := os.Getenv("APP_NAME")
			if appName == "" {
				appName = "ujds"
			}

			l := logger.New().With().Str("app", appName).Logger()

			cfg := config.Config{}

			fi, err := os.Stat(cfgPath)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				l.Fatal().Err(err).Msgf("failed to get %s info", cfgPath)
				return
			}

			if fi != nil {
				if err := cfgloader.LoadFromPath(cfgPath, &cfg, config.Schema); err != nil {
					l.Fatal().Err(err).Msg("failed to load config")
					return
				}
			}

			if err := cfgloader.LoadFromEnv(appName, &cfg); err != nil {
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
			l.Info().Msg("db connection ok")

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

			a := api.New(cfg.API, db, l.With().Str("pkg", "api").Logger())
			s := server.New(cfg.Server, a, l.With().Str("pkg", "server").Logger())

			if err := s.Run(cmd.Context()); errors.Is(err, http.ErrServerClosed) {
				l.Info().Msg("server stopped")
			} else if err != nil {
				l.Error().Err(err).Msg("")
			}
		},
	}

	cmd.Flags().BoolVar(&migUp, "migrate-up", false, "apply database migrations")
	cmd.Flags().BoolVar(&migDown, "migrate-down", false, "revert database migrations")

	cmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "path to the config file")

	return cmd
}
