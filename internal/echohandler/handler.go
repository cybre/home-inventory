package echohandler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type InputBinder[T any] func(echo.Context) (T, error)
type InputValidator[T any] func(T) error
type Handler[T any] func(echo.Context, T) error

type handler[T any] struct {
	inputBinder    InputBinder[T]
	inputValidator InputValidator[T]
	inputHandler   Handler[T]
}

func NewValidateHandler[T any](inputHandler Handler[T], validate *validator.Validate, opts ...Option[T]) echo.HandlerFunc {
	return NewHandler[T](inputHandler, append(opts, WithValidateInputValidator[T](validate))...)
}

func NewHandler[T any](inputHandler Handler[T], opts ...Option[T]) echo.HandlerFunc {
	handler := handler[T]{
		inputBinder:  DefaultInputBinder[T],
		inputHandler: inputHandler,
	}

	for _, opt := range opts {
		opt(&handler)
	}

	return func(c echo.Context) error {
		data, err := handler.inputBinder(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if handler.inputValidator != nil {
			if err := handler.inputValidator(data); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}

		return handler.inputHandler(c, data)
	}
}

type Option[T any] func(*handler[T])

func WithInputBinder[T any](inputBinder InputBinder[T]) Option[T] {
	return func(h *handler[T]) {
		h.inputBinder = inputBinder
	}
}

func WithInputValidator[T any](inputValidator InputValidator[T]) Option[T] {
	return func(h *handler[T]) {
		h.inputValidator = inputValidator
	}
}

func WithValidateInputValidator[T any](validate *validator.Validate) Option[T] {
	return func(h *handler[T]) {
		h.inputValidator = func(data T) error {
			return validate.Struct(data)
		}
	}
}

func DefaultInputBinder[T any](c echo.Context) (T, error) {
	var data T
	if err := c.Bind(&data); err != nil {
		return data, err
	}

	return data, nil
}
