package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/ashep/datapimp/authservice"
	"github.com/ashep/datapimp/dataservice"
	_ "github.com/lib/pq"

	"github.com/ashep/datapimp/config"
	"github.com/ashep/datapimp/mq"
	"github.com/ashep/datapimp/server"
	"github.com/rs/zerolog"
)

func Run(ctx context.Context, cfg config.Config, l zerolog.Logger) error {
	if cfg.DB.DSN == "" {
		return fmt.Errorf("empty db dsn")
	}

	if cfg.MQ.URI == "" {
		return fmt.Errorf("empty mq uri")
	}

	if envAS := os.Getenv("DATAPIMP_ADMIN_SECRET"); envAS != "" {
		cfg.Auth.AdminSecret = envAS
	}
	if cfg.Auth.AdminSecret == "" {
		return errors.New("empty admin secret")
	} else if len(cfg.Auth.AdminSecret) < 8 {
		return errors.New("too short admin secret")
	}

	mqc, err := mq.NewClient("main", cfg.MQ, l.With().Str("pkg", "mq").Logger())
	if err != nil {
		return fmt.Errorf("failed to initialize a message queue client: %w", err)
	}
	defer mqc.Close()
	mqc.Run(ctx)

	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("failed to open a database: %w", err)
	}

	authSvc := authservice.New(db, cfg.Auth, l.With().Str("pkg", "auth").Logger())

	dataSvc := dataservice.New(cfg.Data, db, mqc, l.With().Str("pkg", "data").Logger())
	if err := dataSvc.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize data service: %w", err)
	}

	return server.New(cfg.Server, authSvc, dataSvc, l.With().Str("pkg", "server").Logger()).Run(ctx)
}
