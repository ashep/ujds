package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashep/ujds/cmd/root"
)

func main() {
	time.Local = time.UTC

	cmd := root.New()

	ctx, ctxC := context.WithCancel(context.Background())
	defer ctxC()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sig
		fmt.Printf("%s signal received\n", s)
		ctxC()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
