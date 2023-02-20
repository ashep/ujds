package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashep/datapimp/cmd/root"
)

func main() {
	cmd := root.New()

	ctx, ctxC := context.WithCancel(context.Background())
	defer ctxC()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigs
		fmt.Printf("%s signal received\n", s)
		ctxC()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
