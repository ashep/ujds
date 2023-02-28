package dataservice

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/datapimp/mq"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	cfg Config
	db  *sql.DB
	mqc *mq.Client
	l   zerolog.Logger

	init bool
}

func New(cfg Config, db *sql.DB, mqc *mq.Client, l zerolog.Logger) *Service {
	if cfg.Queues.Push == "" {
		cfg.Queues.Push = "push"
	}

	return &Service{cfg: cfg, db: db, mqc: mqc, l: l}
}

func (s *Service) Init(ctx context.Context) error {
	if s.init {
		return nil
	}

	ch, err := s.mqc.Channel(ctx)
	if err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(s.cfg.Queues.Push, true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare a push queue: %w", err)
	}
	s.l.Debug().Str("name", s.cfg.Queues.Push).Msg("push queue declared")

	if err := ch.QueueBind(s.cfg.Queues.Push, "", "amq.fanout", false, nil); err != nil {
		return fmt.Errorf("failed to bind a push queue: %w", err)
	}
	s.l.Debug().Str("name", s.cfg.Queues.Push).Str("exchange", "amq.fanout").Msg("push queue bound to exchange")

	s.init = true

	return nil
}

func (s *Service) UpsertItem(ctx context.Context, typ, id string, data []byte) (Item, error) {
	if !s.init {
		return Item{}, errors.New("service is not initialized")
	}

	var (
		item Item
		err  error
	)

	if id != "" {
		item, err = s.UpdateItem(ctx, id, data)
	} else {
		item, err = s.CreateItem(ctx, typ, data)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return item, ErrNotFound
	} else if err != nil {
		return item, fmt.Errorf("failed to update item in db: %w", err)
	}

	if err = s.mqc.Publish(ctx, "amq.fanout", typ, data); err != nil {
		return item, fmt.Errorf("failed to send a message to mq: %w", err)
	}

	return item, nil
}

func (s *Service) CreateItem(ctx context.Context, itemType string, data []byte) (Item, error) {
	if !s.init {
		return Item{}, errors.New("service is not initialized")
	}

	var item Item

	tx, err := s.db.Begin()
	if err != nil {
		return item, err
	}

	id := uuid.NewString()

	_, err = tx.ExecContext(ctx, `INSERT INTO item (id, version, type_id) VALUES ($1, $2, $3)`, id, 0, itemType)
	if err != nil {
		return item, err
	}

	if err := tx.Commit(); err != nil {
		return item, err
	}

	return s.GetItem(ctx, id)
}

func (s *Service) UpdateItem(ctx context.Context, id string, data []byte) (Item, error) {
	if !s.init {
		return Item{}, errors.New("service is not initialized")
	}

	item, err := s.GetItem(ctx, id)
	if err != nil {
		return item, err
	}

	s.l.Debug().
		Int("diff", bytes.Compare(item.Data, data)).
		Str("src", string(item.Data)).
		Str("dst", string(data)).
		Msg("compare")

	if bytes.Compare(item.Data, data) == 0 {
		return item, nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return item, err
	}

	_, err = tx.ExecContext(ctx, `UPDATE item SET version=$1 WHERE id=$2`, 0, id)
	if err != nil {
		return item, err
	}

	if err := tx.Commit(); err != nil {
		return item, err
	}

	item.Data = data

	return item, nil
}

func (s *Service) GetItem(ctx context.Context, id string) (Item, error) {
	if !s.init {
		return Item{}, errors.New("service is not initialized")
	}

	var r Item

	q := `SELECT type, version, data, time FROM item WHERE id=$1 LIMIT 1`

	var (
		typ  string
		ver  uint64
		data []byte
		tm   time.Time
	)

	row := s.db.QueryRowContext(ctx, q, id)

	if err := row.Scan(&typ, &ver, &data, &tm); err != nil {
		return r, fmt.Errorf("failed to execute a db query: %w", err)
	}

	r.Id = id
	r.Type = typ
	r.Version = ver
	r.Time = tm
	r.Data = data

	return r, nil
}
