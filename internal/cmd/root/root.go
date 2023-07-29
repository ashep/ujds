package root

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/ashep/go-cfgloader"
	"github.com/ashep/ujds/internal/api"
	"github.com/ashep/ujds/internal/config"
	"github.com/ashep/ujds/internal/logger"
	"github.com/ashep/ujds/internal/migration"
	"github.com/ashep/ujds/internal/server"
)

func New() *cobra.Command {
	var (
		cfgPath string
		migUp   bool
		migDown bool
	)

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			appName := os.Getenv("APP_NAME")
			if appName == "" {
				appName = "ujds"
			}

			l := logger.New().With().Str("app", appName).Logger()

			cfg, err := loadConfig(appName, cfgPath)
			if err != nil {
				l.Fatal().Err(err).Msg("config load failed")
				return
			}

			db, err := dbConn(cmd.Context(), cfg.DB.DSN)
			if err != nil {
				l.Fatal().Err(err).Msg("database connection failed")
				return
			}

			if migUp {
				if err := migration.Up(db); err != nil {
					l.Fatal().Err(err).Msg("failed to apply migrations")
				}
				l.Info().Msg("migrations applied")
				return
			}

			if migDown {
				if err := migration.Down(db); err != nil {
					l.Fatal().Err(err).Msg("failed to revert migrations")
				}
				l.Info().Msg("migrations reverted")
				return
			}

			a := api.New(db, l.With().Str("pkg", "api").Logger())
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

func loadConfig(appName, cfgPath string) (config.Config, error) {
	cfg := config.Config{}

	fi, err := os.Stat(cfgPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, fmt.Errorf("failed to get %s info: %w", cfgPath, err)
	}

	if fi != nil {
		if err := cfgloader.LoadFromPath(cfgPath, &cfg, config.Schema); err != nil {
			return cfg, fmt.Errorf("failed to load from path: %w", err)
		}
	}

	if err := cfgloader.LoadFromEnv(appName, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to load from env: %w", err)
	}

	return cfg, nil
}

func dbConn(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("empty dsn")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}
