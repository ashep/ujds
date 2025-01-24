package testapp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ashep/go-app/buflogwriter"
	"github.com/ashep/go-app/runner"
	"github.com/ashep/ujds/internal/app"
	"github.com/ashep/ujds/internal/server"
	"github.com/ashep/ujds/sdk/client"
	_ "github.com/lib/pq" // it's ok in tests
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	checkPeriod = time.Millisecond * 100
	checkCount  = 50
	dbDSN       = "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
)

type TestApp struct {
	t      *testing.T
	runner *runner.Runner[*app.App, app.Config]
	stop   context.CancelFunc // shut down the app
	done   chan struct{}      // closed when the app stopped
	db     *TestDB
	l      *buflogwriter.BufLogWriter
}

func New(t *testing.T) *TestApp {
	t.Helper()

	db := newDB(t, dbDSN)

	// Get free port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(fmt.Errorf("listen: %w", err))
	}

	addr := lis.Addr().String()
	require.NoError(t, lis.Close())

	cfg := app.Config{
		DB: app.Database{
			DSN: dbDSN,
		},
		Server: server.Config{
			Address:   addr,
			AuthToken: "theAuthToken",
		},
	}

	ta := &TestApp{
		t:    t,
		db:   db,
		done: make(chan struct{}),
		l:    buflogwriter.New(),
	}

	ta.runner = runner.New(app.New, cfg).
		WithLogWriter(ta.l).
		WithDefaultHTPServer().
		WithStopper(func(cf context.CancelFunc) { ta.stop = cf })

	return ta
}

func (ta *TestApp) Start() *TestApp {
	ta.db.Reset()

	go func() {
		ta.runner.Run()
		close(ta.done)
	}()

	netAddr, err := net.ResolveTCPAddr("tcp", ta.runner.AppConfig().Server.Address)
	require.NoError(ta.t, err)

	tk := time.NewTicker(checkPeriod)
	defer tk.Stop()

	started := false
	for i := 0; i < checkCount && !started; i++ {
		<-tk.C

		if _, err := net.DialTCP("tcp", nil, netAddr); err != nil {
			continue
		}

		started = true
	}

	if !started {
		ta.t.Fatalf("app has not started within %s", checkPeriod*checkCount)
	}

	return ta
}

func (ta *TestApp) Stop() {
	ta.stop()

	tk := time.NewTicker(checkPeriod)
	defer tk.Stop()

	defer ta.db.d.Close()

	for i := 0; i < checkCount; i++ {
		select {
		case <-ta.done:
			return
		case <-tk.C:
			continue
		}
	}

	ta.t.Fatalf("app has not stopped within %s", checkPeriod*checkCount)
}

func (ta *TestApp) Client(authToken string) *client.Client {
	if authToken == "" {
		authToken = ta.runner.AppConfig().Server.AuthToken
	}

	return client.New("http://"+ta.runner.AppConfig().Server.Address, authToken, http.DefaultClient)
}

func (ta *TestApp) Logs() string {
	return ta.l.String()
}

func (ta *TestApp) AssertNoLogErrors() {
	assert.NotContains(ta.t, ta.Logs(), `"level":"error"`)
}

func (ta *TestApp) AssertNoLogWarns() {
	assert.NotContains(ta.t, ta.Logs(), `"level":"warn"`)
}

func (ta *TestApp) DB() *TestDB {
	return ta.db
}
