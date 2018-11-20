package actions

import (
	"fmt"
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
		cu := getCurrentUser(c)
		if !cu.IsAdmin() {
			return forbidden(c, "%v tried to access %v", cu, c.Request().RequestURI)
		}
		return next(c)
	}
}

func permissionKeeper(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		cu := getCurrentUser(c)
		if !cu.IsAdmin() {
			ri, ok := c.Value("current_route").(buffalo.RouteInfo)
			if !ok {
				return e500(c, "incorrect route: %v", c.Value("current_route"))
			}
			action := fmt.Sprintf("%v %v", ri.Method, ri.PathName)
			rp := permissionMap[action]
			if rp != "" && !cu.HasRole(rp) {
				return forbidden(c, "%v has no permission %v for %v", cu, rp, action)
			}
		}
		return next(c)
	}
}

func addPermission(action, perm string) {
	permissionMap[action] = perm
}

var permissionMap = map[string]string{}
