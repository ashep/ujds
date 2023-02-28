package mq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Client struct {
	name string
	cfg  Config
	conn *amqp.Connection
	ch   *amqp.Channel
	l    zerolog.Logger

	mux *sync.RWMutex
}

func NewClient(name string, cfg Config, l zerolog.Logger) (*Client, error) {
	if cfg.URI == "" {
		return nil, errors.New("config.uri is empty")
	}

	l = l.With().Str("conn_name", name).Logger()

	return &Client{name: name, cfg: cfg, l: l, mux: &sync.RWMutex{}}, nil
}

func (c *Client) Run(ctx context.Context) {
	var (
		connClose chan *amqp.Error
		chanClose chan *amqp.Error
		err       error
	)

	go func() {
		n := 0
		for {
			select {
			case <-ctx.Done():
				c.l.Debug().Msg("context done")
				return
			default:
			}

			if n < 30 {
				n++
			}

			connClose = make(chan *amqp.Error)
			chanClose = make(chan *amqp.Error)

			c.mux.Lock()
			c.conn, c.ch, err = c.connect()
			c.mux.Unlock()
			if err != nil {
				c.l.Warn().Err(err).Msgf("failed to connect, retrying in %d seconds", n)
				time.Sleep(time.Second * time.Duration(n))
				continue
			}

			n = 0
			c.l.Debug().Msg("connected")

			c.conn.NotifyClose(connClose)
			c.ch.NotifyClose(chanClose)

			select {
			case <-ctx.Done():
				c.l.Debug().Msg("context done")
				return

			case err := <-connClose:
				c.mux.Lock()
				c.conn = nil
				c.mux.Unlock()
				c.l.Warn().Err(err).Msg("connection closed")

			case err := <-chanClose:
				c.mux.Lock()
				c.ch = nil
				c.mux.Unlock()
				c.l.Warn().Err(err).Msg("channel closed")
			}
		}
	}()

	c.l.Debug().Msg("running")
}

func (c *Client) Channel(ctx context.Context) (*amqp.Channel, error) {
	for {
		c.mux.RLock()
		if c.ch != nil {
			c.mux.RUnlock()
			return c.ch, nil
		}
		c.mux.RUnlock()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (c *Client) Publish(ctx context.Context, exchange, key string, body []byte) error {
	ch, err := c.Channel(ctx)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx, exchange, key, false, false, amqp.Publishing{
		Headers:         nil,
		ContentType:     "",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            body,
	})
}

func (c *Client) Close() {
	var err error

	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.ch != nil && !c.ch.IsClosed() {
		if err = c.ch.Close(); err != nil {
			c.l.Fatal().Err(err).Msg("failed to close the channel")
		}
	}

	if c.conn != nil && !c.conn.IsClosed() {
		if err = c.conn.Close(); err != nil {
			c.l.Fatal().Err(err).Msg("failed to close the connection")
		}
	}
}

func (c *Client) connect() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(c.cfg.URI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return conn, ch, nil
}
