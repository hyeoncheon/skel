package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

func authorizeKeeper(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userID := c.Session().Get("user_id"); userID == nil {
			c.Logger().Warn("unauthorized access to ", c.Request().RequestURI)
			c.Flash().Add("danger", t(c, "Login required"))
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		return next(c)
	}
}

func contextMapper(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("uart_url", envy.Get("UART_URL", "/"))
		c.Set("TIME_FORMAT", "2006-01-02T15:04:05Z07:00")
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

func adminKeeper(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user := getCurrentUser(c)
		if !user.IsAdmin() {
			c.Logger().WithField("category", "security").
				Errorf("%v tried to access %v", user, c.Request().RequestURI)
			c.Flash().Add("danger", t(c, "Staff only"))
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		return next(c)
	}
}
