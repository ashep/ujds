package app

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
	"github.com/ashep/ujds/internal/server/indexhandler"
	"github.com/ashep/ujds/internal/server/recordhandler"
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

	ir := indexrepository.New(db, a.l)
	rr := recordrepository.New(db, a.l)

	ih := indexhandler.New(ir, time.Now, a.l)
	rh := recordhandler.New(ir, rr, time.Now, a.l)
	s := server.New(a.cfg.Server, ih, rh, a.l.With().Str("pkg", "server").Logger())

	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server run failed: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("db close failed: %w", err)
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
