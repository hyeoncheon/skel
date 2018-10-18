package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

func LogoutHandler(c buffalo.Context) error {
	destroySession(c)
	c.Flash().Add("success", t(c, "You have been successfully logged out"))
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
