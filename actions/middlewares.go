package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
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

func contextHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userID := c.Session().Get("user_id"); userID != nil {
			if id, ok := userID.(uuid.UUID); ok {
				c.Set("user_id", id)
			}
			c.Set("user_mail", c.Session().Get("user_mail"))
			c.Set("user_name", c.Session().Get("user_name"))
			c.Set("user_icon", c.Session().Get("user_icon"))
			c.Set("user_roles", c.Session().Get("user_roles"))
		}
		return next(c)
	}
}
