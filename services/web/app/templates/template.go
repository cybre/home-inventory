package templates

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/cybre/home-inventory/services/web/app/auth"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/unrolled/render"

	_ "embed"
)

const (
	CustomBaseLayoutKey = "baseLayout"
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

	data := map[string]interface{}{
		"User":     sess.Values[auth.AuthSessionProfileKey],
		"PageData": pageData,
	}

	return t.r.HTML(w, http.StatusOK, name, data)
}
