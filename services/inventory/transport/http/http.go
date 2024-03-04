package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/bnkamalesh/errors"
	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/internal/middleware"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type HouseholdService interface {
	CreateHousehold(context.Context, shared.CreateHouseholdCommandData) error
	UpdateHousehold(context.Context, shared.UpdateHouseholdCommandData) error
	DeleteHousehold(context.Context, shared.DeleteHouseholdCommandData) error

	AddRoom(context.Context, shared.AddRoomCommandData) error
	UpdateRoom(context.Context, shared.UpdateRoomCommandData) error
	DeleteRoom(context.Context, shared.DeleteRoomCommandData) error

	AddItem(context.Context, shared.AddItemCommandData) error
	UpdateItem(context.Context, shared.UpdateItemCommandData) error

	GetUserHouseholds(context.Context, string) ([]shared.UserHousehold, error)
	GetUserHousehold(context.Context, string, string) (shared.UserHousehold, error)

	GetUserHouseholdRoom(context.Context, string, string, string) (shared.UserHouseholdRoom, error)
}

func NewHTTPTransport(ctx context.Context, serverAddress string, householdService HouseholdService) error {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code, message, _ := errors.HTTPStatusCodeMessage(err)

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.String(code, message)
		}
	}

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

	e.Use(middleware.RequestAndCorrelationIDLogging(logger))
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus:       true,
		LogMethod:       true,
		LogURI:          true,
		LogError:        true,
		LogResponseSize: true,
		HandleError:     true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			logger := logging.FromContext(c.Request().Context())

			if v.Error == nil {
				logger.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.Int64("response_size", v.ResponseSize),
				)
			} else {
				logger.LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR",
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

	buildHouseholdRoutes(e, householdService, validate)

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
