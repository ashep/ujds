package testapp

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq" // it's ok in tests
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/app"
	"github.com/ashep/ujds/internal/server"
)

const (
	tcpPort     = 9000
	checkPeriod = time.Millisecond * 100
	checkCount  = 50
)

type TestApp struct {
	app *app.App
	db  *TestDB
	l   zerolog.Logger
	lb  *strings.Builder
}

func New(t *testing.T) *TestApp {
	t.Helper()

	db := newDB(t)
	db.Reset(t)

	lb := &strings.Builder{}
	l := zerolog.New(lb)

	a := app.New(app.Config{
		DB: app.Database{
			DSN: "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable",
		},
		Server: server.Config{
			Address:   ":" + strconv.Itoa(tcpPort),
			AuthToken: "theAuthToken",
		},
	}, l)

	return &TestApp{
		app: a,
		db:  db,
		l:   l,
		lb:  lb,
	}
}

func (a *TestApp) Start(t *testing.T) func() {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		require.NoError(t, a.app.Run(ctx))
		close(done)
	}()

	started := false

	addr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: tcpPort,
	}

	tk1 := time.NewTicker(checkPeriod)
	defer tk1.Stop()

	for i := 0; i < checkCount; i++ {
		<-tk1.C

		if _, err := net.DialTCP("tcp", nil, &addr); err != nil {
			continue
		}

		started = true

		break
	}

	if !started {
		t.Fatalf("app has not started within %s", checkPeriod*checkCount)
	}

	return func() {
		cancel()
		<-done

		tk2 := time.NewTicker(checkPeriod)
		defer tk2.Stop()

		for i := 0; i < checkCount; i++ {
			select {
			case <-done:
				return
			case <-tk2.C:
				continue
			}
		}

		t.Fatalf("app has not stopped within %s", checkPeriod*checkCount)
	}
}

func (a *TestApp) Logs() string {
	return a.lb.String()
}

func (a *TestApp) AssertNoLogErrors(t *testing.T) {
	t.Helper()
	assert.NotContains(t, a.Logs(), `"level":"error"`)
}

func (a *TestApp) AssertNoLogWarns(t *testing.T) {
	t.Helper()
	assert.NotContains(t, a.Logs(), `"level":"warn"`)
}

func (a *TestApp) DB() *TestDB {
	return a.db
}
