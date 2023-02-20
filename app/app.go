package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ashep/datapimp/dataservice"
	_ "github.com/lib/pq"

	"github.com/ashep/datapimp/config"
	"github.com/ashep/datapimp/mq"
	"github.com/ashep/datapimp/server"
	"github.com/rs/zerolog"
)

func Run(ctx context.Context, cfg config.Config, db *sql.DB, l zerolog.Logger) error {
	if cfg.MQ.URI == "" {
		cfg.MQ.URI = "amqp://guest:guest@localhost:5672/datapimp"
	}

	mqc, err := mq.NewClient("main", cfg.MQ, l.With().Str("pkg", "mq").Logger())
	if err != nil {
		return fmt.Errorf("failed to initialize a message queue client: %w", err)
	}
	defer mqc.Close()

	mqc.Run(ctx)

	ds, err := dataservice.New(cfg.Data, db, mqc, l.With().Str("pkg", "data").Logger())
	if err != nil {
		return fmt.Errorf("failed to initialize data service: %w", err)
	}

	return server.New(cfg.Server, ds, l.With().Str("pkg", "server").Logger()).Run(ctx)
}
