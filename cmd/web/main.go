package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/services/web/app"
)

const serviceName = "web"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("service", serviceName)
	slog.SetDefault(logger)

	ctx = logging.WithLogger(ctx, logger)

	auth, err := authenticator.New()
	if err != nil {
		panic(err)
	}

	if err := app.New(ctx, auth, logger); err != nil {
		panic(err)
	}
}
