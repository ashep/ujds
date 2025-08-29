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

	"connectrpc.com/connect"
	"github.com/ashep/go-app/runner"
	"github.com/ashep/ujds/internal/indexrepository"
	"github.com/ashep/ujds/internal/recordrepository"
	"github.com/ashep/ujds/internal/rpc/indexhandler"
	"github.com/ashep/ujds/internal/rpc/recordhandler"
	"github.com/ashep/ujds/internal/validation"
	"github.com/ashep/ujds/migration"
	indexconnect "github.com/ashep/ujds/sdk/proto/ujds/index/v1/v1connect"
	recordconnect "github.com/ashep/ujds/sdk/proto/ujds/record/v1/v1connect"
	"github.com/rs/zerolog"
)

type App struct {
	cfg *Config
	rt  *runner.Runtime
	l   zerolog.Logger
}

func New(cfg *Config, rt *runner.Runtime) (*App, error) {
	return &App{
		cfg: cfg,
		rt:  rt,
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

	interceptors := connect.WithInterceptors(
		auth(a.cfg.Server.AuthToken),
	)

	indexPath, indexHandler := indexconnect.NewIndexServiceHandler(
		indexhandler.New(ir, time.Now, a.l),
		interceptors,
	)
	a.rt.Server.Handle(indexPath, cors(indexHandler))

	recordPath, recordHandler := recordconnect.NewRecordServiceHandler(
		recordhandler.New(ir, rr, time.Now, a.l),
		interceptors,
	)
	a.rt.Server.Handle(recordPath, cors(recordHandler))

	var resErr error
	a.rt.Logger.Info().Str("addr", a.rt.Server.Listener().Addr().String()).Msg("starting server")
	if srvErr := <-a.rt.Server.Start(ctx); srvErr != nil {
		if !errors.Is(srvErr, http.ErrServerClosed) {
			resErr = errors.Join(resErr, fmt.Errorf("server: %w", err))
		}
	}
	if dbCloseErr := db.Close(); dbCloseErr != nil {
		resErr = errors.Join(resErr, fmt.Errorf("db close: %w", dbCloseErr))
	}

	return resErr
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

func auth(token string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if token == "" {
				return next(ctx, req)
			}

			if token == strings.ReplaceAll(req.Header().Get("Authorization"), "Bearer ", "") {
				return next(ctx, req)
			}

			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authorized"))
		}
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
