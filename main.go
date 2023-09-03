package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/ashep/go-apprun"

	"github.com/ashep/ujds/internal/app"
)

func main() {
	apprun.Run("ujds", app.New, app.Config{})
}
