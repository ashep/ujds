package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/ashep/ujds/internal/cmd/root"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo // it's ok here
		os.Exit(1)
	}
}

func run() error {
	time.Local = time.UTC

	cmd := root.New()

	ctx, ctxC := context.WithCancel(context.Background())
	defer ctxC()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sig
		fmt.Printf("%s signal received\n", s) //nolint:forbidigo // it's ok
		ctxC()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("execute failed: %w", err)
	}

	return nil
}
