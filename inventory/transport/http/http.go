package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/cybre/home-inventory/inventory/shared"
	"github.com/cybre/home-inventory/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HouseholdService interface {
	CreateHousehold(context.Context, shared.CreateHouseholdCommandData) error
	AddRoom(context.Context, shared.AddRoomCommandData) error
	AddItem(context.Context, shared.AddItemCommandData) error
	UpdateItem(context.Context, shared.UpdateItemCommandData) error
}

type UserService interface {
	GenerateLoginToken(context.Context, shared.GenerateLoginTokenCommandData) error
	CreateUser(context.Context, shared.CreateUserCommandData) error
	LoginViaToken(context.Context, shared.LoginViaTokenCommandData) error
}

func NewHTTPTransport(ctx context.Context, householdService HouseholdService, userService UserService) error {
	e := echo.New()

	logger := logging.FromContext(ctx)

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		for _, tag := range []string{"json", "form", "query", "param"} {
			name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
			if name == "-" {
				return ""
			}

			if name != "" {
				return name
			}
		}

		return fld.Name
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			loggerCtx := logging.WithLogger(ctx, logger.With(slog.String("correlation_id", uuid.NewString())))
			c.SetRequest(c.Request().WithContext(loggerCtx))

			return next(c)
		}
	})

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
					slog.String("request_id", v.RequestID),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	buildHouseholdRoutes(e, householdService, validate)
	buildUserRoutes(e, userService, validate)

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
