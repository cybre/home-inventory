package http

import (
	"net/http"

	"github.com/cybre/home-inventory/inventory/shared"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func buildUserRoutes(e *echo.Echo, userService UserService, validate *validator.Validate) {
	e.POST("/user/login", validatedHandler(loginHandler(userService), validate))
	e.POST("/user/register", validatedHandler(registerHandler(userService), validate))
}

func loginHandler(userService UserService) InputHandler[shared.GenerateOneTimeTokenCommandData] {
	return func(c echo.Context, data shared.GenerateOneTimeTokenCommandData) error {
		if err := userService.GenerateOneTimeToken(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func registerHandler(userService UserService) InputHandler[shared.CreateUserCommandData] {
	return func(c echo.Context, data shared.CreateUserCommandData) error {
		if err := userService.CreateUser(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}
