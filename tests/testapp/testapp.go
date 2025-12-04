package testapp

import (
	"net/http"
	"testing"
	"time"

	"github.com/ashep/go-app/testlogger"
	"github.com/ashep/go-app/testrunner"
	"github.com/ashep/ujds/internal/app"
	"github.com/ashep/ujds/sdk/client"
	_ "github.com/lib/pq" // it's ok in tests
)

type tRunner interface {
	Logger() *testlogger.Logger
}

type TestApp struct {
	t   *testing.T
	cfg app.Config
	rnr tRunner
	db  *TestDB
}

func New(t *testing.T) *TestApp {
	t.Helper()

	db := newDB(t)
	cfg := app.Config{
		DB: app.Database{
			DSN: db.DSN,
		},
		Server: app.Server{
			Addr:      testrunner.RandLocalTCPAddr(t).String(),
			AuthToken: "theAuthToken",
		},
	}

	rnr := testrunner.New(t, app.Run, cfg).
		SetHTTPReadyStartWaiter("http://"+cfg.Server.Addr+"/metrics", time.Second*5).
		Start()

	ta := &TestApp{
		t:   t,
		cfg: cfg,
		rnr: rnr,
		db:  db,
	}

	return ta
}

func (ta *TestApp) Client(authToken string) *client.Client {
	if authToken == "" {
		authToken = ta.cfg.Server.AuthToken
	}
	return client.New("http://"+ta.cfg.Server.Addr, authToken, http.DefaultClient)
}

func (ta *TestApp) AssertNoWarnsAndErrors() {
	ta.rnr.Logger().AssertNoWarnsAndErrors()
}

func (ta *TestApp) DB() *TestDB {
	return ta.db
}
