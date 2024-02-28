package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cybre/home-inventory/internal/authenticator"
	internalcache "github.com/cybre/home-inventory/internal/cache"
	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/internal/middleware"
	inventoryclient "github.com/cybre/home-inventory/services/inventory/client"
	"github.com/cybre/home-inventory/services/web/app/routes"
	"github.com/cybre/home-inventory/services/web/app/templates"
	"github.com/eko/gocache/lib/v4/cache"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, serverAddress string, logger *slog.Logger) error {
	e := echo.New()
	e.Renderer = templates.New()

	authenticator, err := authenticator.New()
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_CLIENT_CACHE_ADDRRESS"),
	}))
	cacheManager := cache.New[string](redisStore)
	cache := internalcache.New(cacheManager, 2*time.Minute)
	inventoryClient := inventoryclient.New(os.Getenv("INVENTORY_API"), cache)

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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(routes.LoadHouseholdsIntoContext(inventoryClient))

	e.Static("/static", "static")

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
