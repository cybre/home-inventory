package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/services/web/app"
	"github.com/joho/godotenv"
)

const serviceName = "web"

var serverAddress = os.Getenv("SERVER_ADDRESS")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("service", serviceName)
	slog.SetDefault(logger)

	ctx = logging.WithLogger(ctx, logger)

	if os.Getenv("ENV") == "dev" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}

		serverAddress = os.Getenv("SERVER_ADDRESS")
	}

	if err := app.New(ctx, serverAddress, logger); err != nil {
		panic(err)
	}
}
