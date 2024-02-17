package app

import (
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/cybre/home-inventory/internal/echomiddleware"
	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/services/web/app/callback"
	"github.com/cybre/home-inventory/services/web/app/home"
	"github.com/cybre/home-inventory/services/web/app/login"
	"github.com/cybre/home-inventory/services/web/app/logout"
	"github.com/cybre/home-inventory/services/web/app/postlogout"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(ctx context.Context, auth *authenticator.Authenticator, logger *slog.Logger) error {
	e := echo.New()

	gob.Register(map[string]interface{}{})

	e.Use(echomiddleware.RequestAndCorrelationIDLogging(logger))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogRemoteIP: true,
		LogMethod:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
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

	e.Use(middleware.Recover())

	// TODO: Load key from env
	store := sessions.NewCookieStore([]byte("secret"))
	e.Use(session.Middleware(store))

	e.GET("/login", login.Handler(auth))
	e.GET("/callback", callback.Handler(auth))
	e.GET("/logout", logout.Handler())
	e.GET("/postlogout", postlogout.Handler())
	e.GET("/", home.Handler())
	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "protected")
	}, isAuthenticated)

	go func() {
		if err := e.Start(":8080"); err != nil {
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

func isAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("auth-session", c)
		if err != nil {
			return err
		}

		if sess.Values["profile"] == nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		logger := logging.FromContext(c.Request().Context())
		logger.Info("user is authenticated", slog.Any("profile", sess.Values["profile"]))

		return next(c)
	}
}
