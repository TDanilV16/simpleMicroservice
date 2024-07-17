package main

import (
	"context"
	"os"
	"os/signal"
	"simpleMicroservice/internal/config"
	"simpleMicroservice/pkg/logger"
	"time"
)

func main() {
	ctx := context.Background()

	ctx = logger.ContextWithTimeStamp(ctx, time.Now())

	cfg := config.MustLoad()

	ctx = logger.ContextWithAddressString(ctx, cfg.Address)
	ctx = logger.ContextWithEnvString(ctx, cfg.Env)

	logger.SetUpLogger(cfg)

	logger.Info(ctx, "initializing server")

	if err := run(ctx); err != nil {
		logger.Fatalf(ctx, "err: %v", err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	<-ctx.Done()
	return ctx.Err()
}
