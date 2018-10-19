package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// ProfileShow handles current user's service profile
func ProfileShow(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("profile.html"))
}
