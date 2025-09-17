package testapp

import (
	"net/http"
	"testing"

	"github.com/ashep/go-app/testrunner"
	"github.com/ashep/ujds/internal/app"
	"github.com/ashep/ujds/sdk/client"
	_ "github.com/lib/pq" // it's ok in tests
)

const (
	dbDSN = "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
)

type tRunner interface {
	Logs() string
	AssertLogNoErrors()
	AssertLogNoWarns()
}

type TestApp struct {
	t   *testing.T
	cfg app.Config
	rnr tRunner
	db  *TestDB
}

func New(t *testing.T) *TestApp {
	t.Helper()

	db := newDB(t, dbDSN).Reset()
	srvAddr := testrunner.RandLocalTCPAddr(t)

	cfg := app.Config{
		DB: app.Database{
			DSN: dbDSN,
		},
		Server: app.Server{
			Addr:      srvAddr.String(),
			AuthToken: "theAuthToken",
		},
	}

	rnr := testrunner.New(t, app.Run, cfg).
		SetTCPPortStartWaiter(srvAddr).Start()

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

func (ta *TestApp) Logs() string {
	return ta.rnr.Logs()
}

func (ta *TestApp) AssertLogNoErrors() {
	ta.rnr.AssertLogNoErrors()
}

func (ta *TestApp) AssertLogNoWarns() {
	ta.rnr.AssertLogNoWarns()
}

func (ta *TestApp) DB() *TestDB {
	return ta.db
}
