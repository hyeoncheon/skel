package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func ProfileShow(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("profile.html"))
}
