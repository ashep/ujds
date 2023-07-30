package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/api"
	"github.com/ashep/ujds/internal/server"
)

type App struct {
	cfg Config
	l   zerolog.Logger
}

func New(cfg Config, l zerolog.Logger) *App {
	return &App{cfg: cfg, l: l}
}

func (a *App) Run(ctx context.Context) error {
	db, err := dbConn(ctx, a.cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	s := server.New(
		a.cfg.Server,
		api.New(db, a.l.With().Str("pkg", "api").Logger()),
		a.l.With().Str("pkg", "server").Logger(),
	)

	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}

func dbConn(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
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
