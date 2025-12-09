package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/go-app/dbmigrator"
	"github.com/ashep/go-app/health"
	"github.com/ashep/go-app/httpserver"
	"github.com/ashep/go-app/prommetrics"
	"github.com/ashep/go-app/runner"
	"github.com/ashep/ujds/internal/indexrepo"
	"github.com/ashep/ujds/internal/recordrepo"
	"github.com/ashep/ujds/internal/rpc/indexhandler"
	"github.com/ashep/ujds/internal/rpc/recordhandler"
	"github.com/ashep/ujds/internal/validation"
	indexconnect "github.com/ashep/ujds/sdk/proto/ujds/index/v1/v1connect"
	recordconnect "github.com/ashep/ujds/sdk/proto/ujds/record/v1/v1connect"
	"github.com/ashep/ujds/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func Run(rt *runner.Runtime[Config]) error { //nolint:cyclop // to do
	l := rt.Log
	cfg := rt.Cfg

	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":9000"
	}

	migRes, err := dbmigrator.RunPostgres(cfg.DB.DSN, sql.FS, "migrations", l)
	if err != nil {
		return fmt.Errorf("migrate db: %w", err)
	}
	if migRes.PrevVersion != migRes.NewVersion {
		l.Info().
			Uint("from", migRes.PrevVersion).
			Uint("to", migRes.NewVersion).
			Msg("database migrated")
	}

	pgx, err := pgxpool.New(rt.Ctx, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}
	defer pgx.Close()

	db := stdlib.OpenDBFromPool(pgx)

	srv := httpserver.New(httpserver.WithAddr(rt.Cfg.Server.Addr))
	health.RegisterServer(srv)
	prommetrics.RegisterServer(rt.AppName, rt.AppVersion, srv)

	idxNameValidator := validation.NewIndexNameValidator()
	recIDValidator := validation.NewRecordIDValidator()
	recDataValidator := validation.NewJSONValidator(cfg.Validation.IndexStruct)

	ir := indexrepo.New(db, idxNameValidator, rt.Log)
	rr := recordrepo.New(db, idxNameValidator, recIDValidator, rt.Log)

	icps := connect.WithInterceptors(auth(rt.Cfg.Server.AuthToken))

	indexPath, indexHandler := indexconnect.NewIndexServiceHandler(
		indexhandler.New(ir, time.Now, rt.Log),
		icps,
	)
	srv.Handle(indexPath, cors(indexHandler))

	recordPath, recordHandler := recordconnect.NewRecordServiceHandler(
		recordhandler.New(ir, rr, idxNameValidator, recIDValidator, recDataValidator, time.Now, rt.Log),
		icps,
	)
	srv.Handle(recordPath, cors(recordHandler))

	var resErr error
	rt.Log.Info().Str("addr", srv.Listener().Addr().String()).Msg("starting")
	if srvErr := srv.Run(rt.Ctx); srvErr != nil {
		if !errors.Is(srvErr, http.ErrServerClosed) {
			resErr = errors.Join(resErr, fmt.Errorf("server: %w", srvErr))
		}
	}

	return resErr
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
