package logout

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

func Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
		if err != nil {
			return fmt.Errorf("failed to build logout URL: %w", err)
		}

		scheme := "http"
		if c.Request().TLS != nil {
			scheme = "https"
		}

		returnTo, err := url.Parse(scheme + "://" + c.Request().Host + "/postlogout")
		if err != nil {
			return fmt.Errorf("failed to build returnTo URL: %w", err)
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
		logoutUrl.RawQuery = parameters.Encode()

		return c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	}
}
