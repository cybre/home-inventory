package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/internal/middleware"
	inventoryclient "github.com/cybre/home-inventory/services/inventory/client"
	"github.com/cybre/home-inventory/services/web/app/routes"
	"github.com/cybre/home-inventory/services/web/app/templates"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func New(ctx context.Context, serverAddress string, logger *slog.Logger) error {
	e := echo.New()
	e.Renderer = templates.New()

	authenticator, err := authenticator.New()
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	e.Use(middleware.RequestAndCorrelationIDLogging(logger))
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus:   true,
		LogRemoteIP: true,
		LogMethod:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			logger := logging.FromContext(c.Request().Context())

			if v.Error == nil {
				logger.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
					slog.String("remote_ip", v.RemoteIP),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR",
					slog.String("remote_ip", v.RemoteIP),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}

			return nil
		},
	}))
	e.Use(echomiddleware.Recover())

	// TODO: Load key from env
	store := sessions.NewCookieStore([]byte("secret"))
	e.Use(session.Middleware(store))

	e.Static("/static", "static")

	inventoryClient := inventoryclient.New(os.Getenv("INVENTORY_API"))

	routes.Initialize(e, authenticator, inventoryClient)

	go func() {
		if err := e.Start(serverAddress); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			panic(err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
