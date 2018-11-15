package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

type intError struct {
	code    string
	message string
}

// error codes. code format is ES + DOmain + Function Code + Error Number
var (
	ESA0SU01 = intError{"ESA0SU01", "no transaction found"}
	ESA0SU02 = intError{"ESA0SU02", "could not store user"}
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
