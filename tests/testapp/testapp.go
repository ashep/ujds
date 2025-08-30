package testapp

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ashep/go-app/buflogwriter"
	"github.com/ashep/go-app/httpserver"
	"github.com/ashep/go-app/runner"
	"github.com/ashep/ujds/internal/app"
	"github.com/ashep/ujds/sdk/client"
	_ "github.com/lib/pq" // it's ok in tests
	"github.com/stretchr/testify/assert"
)

const (
	checkPeriod = time.Millisecond * 100
	checkCount  = 50
	dbDSN       = "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
)

type TestApp struct {
	t      *testing.T
	cfg    *app.Config
	runner *runner.Runner[*app.App, app.Config]
	srv    *httpserver.Server
	stop   context.CancelFunc // shut down the app
	done   chan struct{}      // closed when the app stopped
	db     *TestDB
	l      *buflogwriter.BufLogWriter
}

func New(t *testing.T) *TestApp {
	t.Helper()

	db := newDB(t, dbDSN)

	cfg := &app.Config{
		DB: app.Database{
			DSN: dbDSN,
		},
		Server: app.Server{
			AuthToken: "theAuthToken",
		},
	}

	ta := &TestApp{
		t:    t,
		cfg:  cfg,
		db:   db,
		srv:  httpserver.New(httpserver.WithRandomLocalAddr()),
		done: make(chan struct{}),
		l:    buflogwriter.New(),
	}

	ta.runner = runner.New(app.New).
		WithConfig(cfg).
		WithLogWriter(ta.l).
		WithHTTPServer(ta.srv)

	return ta
}

func (ta *TestApp) Start() *TestApp {
	ta.db.Reset()

	ctx, ctxC := context.WithCancel(context.Background())
	ta.stop = ctxC

	go func() {
		ta.runner.RunContext(ctx)
		close(ta.done)
	}()

	tk := time.NewTicker(checkPeriod)
	defer tk.Stop()

	// Wait the app is up and running
	srvAddr := ta.srv.Listener().Addr().(*net.TCPAddr)
	started := false
	for i := 0; i < checkCount && !started; i++ {
		<-tk.C
		if _, err := net.DialTCP("tcp", nil, srvAddr); err != nil {
			continue
		}
		started = true
	}
	if !started {
		ta.t.Fatalf("app has not started within %s", checkPeriod*checkCount)
	}

	ta.t.Cleanup(ta.shutdown)

	return ta
}

func (ta *TestApp) shutdown() {
	ta.stop()

	tk := time.NewTicker(checkPeriod)
	defer tk.Stop()

	defer ta.db.d.Close() // nolint:errcheck // it's ok in tests

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
		authToken = ta.cfg.Server.AuthToken
	}

	return client.New("http://"+ta.srv.Listener().Addr().String(), authToken, http.DefaultClient)
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
