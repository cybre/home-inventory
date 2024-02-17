package echomiddleware

import (
	"log/slog"

	"github.com/cybre/home-inventory/internal/logging"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestAndCorrelationIDLogging(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Header.Get("X-Request-ID")
			if requestId == "" {
				requestId = uuid.NewString()
				c.Request().Header.Set("X-Request-ID", requestId)
			}

			correlationId := c.Request().Header.Get("X-Correlation-ID")
			if correlationId != "" {
				correlationId = uuid.NewString()
				c.Request().Header.Set("X-Correlation-ID", correlationId)
			}

			loggerCtx := logging.WithLogger(c.Request().Context(), logger.With(slog.String("request_id", requestId), slog.String("correlation_id", correlationId)))
			c.SetRequest(c.Request().WithContext(loggerCtx))

			return next(c)
		}
	}
}
