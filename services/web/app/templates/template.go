package templates

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/cybre/home-inventory/services/web/app/auth"
	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/cybre/home-inventory/services/web/app/htmx"
	"github.com/cybre/home-inventory/services/web/app/routes"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/unrolled/render"

	_ "embed"
)

//go:embed views/*.html
var embeddedTemplates embed.FS

type Renderer struct {
	r *render.Render
}

func New() *Renderer {
	r := render.New(render.Options{
		Layout:    "layout",
		Directory: "views",
		FileSystem: &render.EmbedFileSystem{
			FS: embeddedTemplates,
		},
		Extensions: []string{".html"},
		Funcs: []template.FuncMap{
			{
				"dict": func(values ...interface{}) (map[string]interface{}, error) {
					if len(values)%2 != 0 {
						return nil, errors.New("invalid dict call")
					}
					dict := make(map[string]interface{}, len(values)/2)
					for i := 0; i < len(values); i += 2 {
						key, ok := values[i].(string)
						if !ok {
							return nil, errors.New("dict keys must be strings")
						}
						dict[key] = values[i+1]
					}
					return dict, nil
				},
			},
		},
	})

	return &Renderer{
		r: r,
	}
}

func (t *Renderer) Render(w io.Writer, name string, pageData interface{}, c echo.Context) error {
	sess, err := session.Get(auth.AuthSessionCookieName, c)
	if err != nil {
		return fmt.Errorf("failed to get auth session: %w", err)
	}

	households, ok := helpers.ContextGet[[]shared.UserHousehold](c, routes.ContextHouseholdsKey)
	if !ok {
		households = []shared.UserHousehold{}
	}

	opts := []render.HTMLOptions{}
	if htmx.IsHTMXRequest(c) {
		opts = append(opts, render.HTMLOptions{Layout: "empty_layout"})
	} else {
		pageData = map[string]interface{}{
			"ShouldShowSidebar": len(households) > 0,
			"User":              sess.Values[auth.AuthSessionProfileKey],
			"Households":        households,
			"PageData":          pageData,
		}
	}

	fmt.Printf("opts: %v\n", opts)

	return t.r.HTML(w, http.StatusOK, name, pageData, opts...)
}
