package http

import (
	"net/http"

	"github.com/cybre/home-inventory/inventory/shared"
	eh "github.com/cybre/home-inventory/pkg/echohandler"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func buildUserRoutes(e *echo.Echo, userService UserService, validate *validator.Validate) {
	e.POST("/user/login", eh.NewValidateHandler(loginHandler(userService), validate))
	e.GET("/user/login/:token", eh.NewValidateHandler(loginViaTokenHandler(userService), validate, eh.WithInputBinder(loginViaTokenBinder)))
	e.POST("/user/register", eh.NewValidateHandler(registerHandler(userService), validate))
}

func loginHandler(userService UserService) eh.Handler[shared.GenerateLoginTokenCommandData] {
	return func(c echo.Context, data shared.GenerateLoginTokenCommandData) error {
		if err := userService.GenerateLoginToken(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func loginViaTokenHandler(userService UserService) eh.Handler[shared.LoginViaTokenCommandData] {
	return func(c echo.Context, data shared.LoginViaTokenCommandData) error {
		if err := userService.LoginViaToken(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func loginViaTokenBinder(c echo.Context) (shared.LoginViaTokenCommandData, error) {
	return shared.LoginViaTokenCommandData{
		Token:     c.Param("token"),
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
	}, nil
}

func registerHandler(userService UserService) eh.Handler[shared.CreateUserCommandData] {
	return func(c echo.Context, data shared.CreateUserCommandData) error {
		if err := userService.CreateUser(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}
