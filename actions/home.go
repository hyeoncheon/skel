package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up a home page.
func HomeHandler(c buffalo.Context) error {
	if getCurrentUser(c).IsValid() {
		return c.Redirect(http.StatusTemporaryRedirect, "profilePath()")
	}
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

func LogoutHandler(c buffalo.Context) error {
	destroySession(c)
	c.Flash().Add("success", t(c, "You have been successfully logged out"))
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
