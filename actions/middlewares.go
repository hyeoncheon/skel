package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func AuthorizeHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userID := c.Session().Get("user_id"); userID == nil {
			c.Logger().Warn("unauthorized access to ", c.Request().RequestURI)
			c.Flash().Add("danger", t(c, "Login required"))
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		return next(c)
	}
}

func contextHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userID := c.Session().Get("user_id"); userID != nil {
			user := getCurrentUser(c)
			c.Set("current_user", user)
			c.Set("user_id", user.ID)
			c.Set("user_mail", user.Email)
			c.Set("user_name", user.Name)
			c.Set("user_icon", user.Avatar)
			c.Set("user_roles", user.Roles)
		}
		return next(c)
	}
}
