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
}

func New(cfg Config, db *sql.DB, mcq *mq.Client, l zerolog.Logger) (*Service, error) {
	if cfg.Queues.Push == "" {
		cfg.Queues.Push = "push"
	}

	ch, err := mcq.Channel(context.Background())
	if err != nil {
		return nil, err
	}

	if _, err := ch.QueueDeclare(cfg.Queues.Push, true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("failed to declare a push queue: %w", err)
	}
	l.Debug().Str("name", cfg.Queues.Push).Msg("push queue declared")

	if err := ch.QueueBind(cfg.Queues.Push, "", "amq.fanout", false, nil); err != nil {
		return nil, fmt.Errorf("failed to bind a push queue: %w", err)
	}
	l.Debug().Str("name", cfg.Queues.Push).Str("exchange", "amq.fanout").Msg("push queue bound to exchange")

	return &Service{cfg: cfg, db: db, mqc: mcq, l: l}, nil
}

func (s *Service) Push(ctx context.Context, typ, id string, data []byte) (Item, error) {
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
	var item Item

	tx, err := s.db.Begin()
	if err != nil {
		return item, err
	}

	id := uuid.NewString()

	ver, err := s.insertItemHistory(ctx, tx, id, data)
	if err != nil {
		return item, err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO item (id, version, type) VALUES ($1, $2, $3)`, id, ver, itemType)
	if err != nil {
		return item, err
	}

	if err := tx.Commit(); err != nil {
		return item, err
	}

	return s.GetItem(ctx, id)
}

func (s *Service) UpdateItem(ctx context.Context, id string, data []byte) (Item, error) {
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

	ver, err := s.insertItemHistory(ctx, tx, id, data)
	if err != nil {
		return item, err
	}

	_, err = tx.ExecContext(ctx, `UPDATE item SET version=$1 WHERE id=$2`, ver, id)
	if err != nil {
		return item, err
	}

	if err := tx.Commit(); err != nil {
		return item, err
	}

	item.Version = ver

	return item, nil
}

func (s *Service) GetItem(ctx context.Context, id string) (Item, error) {
	var r Item

	q := `SELECT i.type, i.version, iv.data, iv.time FROM item i 
    LEFT JOIN item_version iv ON i.id = iv.item_id 
    WHERE i.id=$1 LIMIT 1`

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

func (s *Service) insertItemHistory(ctx context.Context, tx *sql.Tx, id string, data []byte) (uint64, error) {
	var ver uint64

	row := tx.QueryRowContext(ctx, `INSERT INTO item_version (item_id, data) VALUES ($1, $2) RETURNING id`, id, data)

	if err := row.Scan(&ver); err != nil {
		return ver, fmt.Errorf("failed to scan version column value: %w", err)
	}

	return ver, nil
}
