package main

import (
	"github.com/ashep/go-app/runner"
	"github.com/ashep/ujds/internal/app"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	appName = "" //nolint:gochecknoglobals // set externally
	appVer  = "" //nolint:gochecknoglobals // set externally
)

func main() {
	runner.New(app.New, app.Config{}).
		WithExtConfig().
		WithConsoleLogWriter().
		WithDefaultHTTPLogWriter(false).
		WithDefaultHTPServer().
		WithMetricsHandler().
		Run()
}
