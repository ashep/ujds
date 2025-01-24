package app

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ashep/go-app/runner"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/indexrepository"
	"github.com/ashep/ujds/internal/recordrepository"
	"github.com/ashep/ujds/internal/server"
	"github.com/ashep/ujds/internal/server/indexhandler"
	"github.com/ashep/ujds/internal/server/recordhandler"
	"github.com/ashep/ujds/internal/validation"
	"github.com/ashep/ujds/migration"
)

type App struct {
	cfg Config
	l   zerolog.Logger
}

func New(cfg Config, rt *runner.Runtime) (*App, error) {
	return &App{
		cfg: cfg,
		l:   rt.Logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error { //nolint:cyclop // to do
	args := []string{os.Args[0]}

	for _, v := range os.Args[1:] {
		if !strings.HasPrefix(v, "-test.") {
			args = append(args, v)
		}
	}

	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	migUp := flagSet.Bool("migrate-up", false, "apply database migrations")
	migDown := flagSet.Bool("migrate-down", false, "revert database migrations")

	if err := flagSet.Parse(args[1:]); err != nil {
		return fmt.Errorf("command line arguments parse failed: %w", err)
	}

	db, err := dbConn(ctx, a.cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	if *migUp {
		if err := migration.Up(db); err != nil {
			return fmt.Errorf("migrration apply failed: %w", err)
		}

		a.l.Info().Msg("migrations applied")

		return nil
	}

	if *migDown {
		if err := migration.Down(db); err != nil {
			return fmt.Errorf("migrration revert failed: %w", err)
		}

		a.l.Info().Msg("migrations reverted")

		return nil
	}

	ir := indexrepository.New(db, validation.NewIndexNameValidator(), a.l)
	rr := recordrepository.New(
		db,
		validation.NewIndexNameValidator(),
		validation.NewRecordIDValidator(),
		validation.NewJSONValidator(),
		a.l,
	)

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
