package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func AuthorizeHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userID := c.Session().Get("user_id"); userID == nil {
			c.Logger().Warn("unauthorized access to ", c.Request().RequestURI)
			c.Flash().Add("danger", "login.required")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		return next(c)
	}
}
