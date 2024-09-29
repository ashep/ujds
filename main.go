package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/ashep/go-apprun/apprun"

	"github.com/ashep/ujds/internal/app"
)

var (
	appName = "" //nolint:gochecknoglobals // set externally
	appVer  = "" //nolint:gochecknoglobals // set externally
)

func main() {
	apprun.Run(app.New, app.Config{}, appName, appVer, nil)
}
