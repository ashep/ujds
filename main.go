package main

import (
	"fmt"
	"os"

	"github.com/ashep/go-app/runner"
	"github.com/ashep/ujds/internal/app"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	res := runner.New(app.Run).
		AddConsoleLogWriter().
		LoadEnvConfig().
		LoadConfigFile("config.yml").
		AddHTTPLogWriter().
		Run()

	if res != nil {
		fmt.Println(res.Error())
		os.Exit(1)
	}
}
