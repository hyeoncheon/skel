package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

type intError struct {
	code    string
	message string
}

// error codes. code format is ES + DOmain + Function Code + Error Number
var (
	ESTX0001 = intError{"ESTX0001", "no transaction found"}
	ESA0SU01 = intError{"ESA0SU01", "could not store user"}
	ESU0PS01 = intError{"ESU0PS01", "could not found the user from datastore"}
)

func oops(c buffalo.Context, x intError, e error) error {
	if e != nil {
		c.Logger().Errorf("system error %v: %v (%v)", x.code, x.message, e)
	} else {
		c.Logger().Errorf("system error %v: %v", x.code, x.message)
	}
	return errors.Errorf("Oops! something went wrong! Error Code: %v", x.code)
}

func eeps(c buffalo.Context, message string) error {
	c.Logger().Errorf("service error: %v", message)
	return errors.Errorf("What's going on? %v", message)
}

//** error handling with render
func forbidden(c buffalo.Context, f string, args ...interface{}) error {
	c.Logger().WithField("category", "security").
		Errorf("SECURITY! "+f, args...)
	return c.Render(http.StatusForbidden, r.HTML("errors/violation"))
}

func e500(c buffalo.Context, f string, args ...interface{}) error {
	impossible(c, f, args...)
	return c.Render(http.StatusInternalServerError, r.HTML("errors/500"))
}

//** internal use only
func impossible(c buffalo.Context, f string, args ...interface{}) {
	c.Logger().WithField("category", "bug").
		Errorf("IMPOSSIBLE! "+f, args...)
}
