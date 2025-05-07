package main

import (
	"context"
	"log"

	"github.com/proof-of-work/internal/build"
	"github.com/proof-of-work/internal/config"
	"github.com/proof-of-work/pkg/logger"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	l, err := logger.New(cfg.App.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = l.WithContext(ctx)

	b := build.New(l, cfg)

	go func() {
		b.WaitShutdown(ctx)
		cancel()
	}()

	go func() {
		defer cancel()

		if err := b.RunTCPClient(ctx); err != nil {
			log.Fatalf("failed to run TCP client: %v", err)
		}
	}()

	<-ctx.Done()
}
