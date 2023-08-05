package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/indexrepository"
	"github.com/ashep/ujds/internal/recordrepository"
	"github.com/ashep/ujds/internal/server"
	"github.com/ashep/ujds/internal/server/handler"
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

	h := handler.New(indexrepository.New(db, a.l), recordrepository.New(db, a.l), time.Now, a.l)
	s := server.New(a.cfg.Server, h, h, a.l.With().Str("pkg", "server").Logger())

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
