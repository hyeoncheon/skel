package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/cloudfoundry"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/skel/models"
)

func init() {
	gothic.Store = App().SessionStore

	uartProvider := cloudfoundry.New(
		os.Getenv("UART_URL"),
		os.Getenv("UART_KEY"),
		os.Getenv("UART_SECRET"),
		fmt.Sprintf("%s%s", App().Host, "/auth/uart/callback"),
		"profile",
	)
	uartProvider.SetName("uart")

	goth.UseProviders(
		uartProvider,
	)
}

// AuthCallback handles callback url of uart OAuth2 and authorize user login
func AuthCallback(c buffalo.Context) error {
	oau, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}
	c.Logger().Infof("attempt to login: %v/%v", oau.UserID, oau.Email)

	if err := validateUARTUser(&oau); err != nil {
		joau, _ := json.Marshal(oau)
		c.Logger().Warnf("invalid uart user: %v (%v)", err, string(joau))
		c.Flash().Add("danger", t(c, err.Error()))
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	user, err := setUser(c, &oau)
	if err != nil {
		c.Flash().Add("danger", t(c, err.Error()))
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	makeSession(c, user)

	c.Logger().Infof("user %v logged in", user)

	c.Flash().Add("success", t(c, "Welcome back! I missed you"))
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func validateUARTUser(oau *goth.User) error {
	roles, ok := oau.RawData["roles"].([]interface{})
	if !ok || len(roles) < 1 {
		return errors.New("You have no access permission")
	}
	return nil
}

// setUser creates or updates user in database, and returns user.
func setUser(c buffalo.Context, oau *goth.User) (*models.User, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return &models.User{}, oops(c, ESTX0001, nil)
	}

	var verrs *validate.Errors
	user := &models.User{Email: oau.Email}
	err := tx.Find(user, oau.UserID)
	if err == nil {
		verrs, err = tx.ValidateAndUpdate(user)
	} else {
		user.ID, _ = uuid.FromString(oau.UserID)
		verrs, err = tx.ValidateAndCreate(user)
	}
	if err != nil {
		return &models.User{}, oops(c, ESA0SU01, err)
	}
	if verrs.HasAny() {
		return &models.User{}, eeps(c, "Validation failed")
	}

	// name, avatar icon, and roles are not stored on database.
	user.Name = oau.Name
	user.Avatar = oau.RawData["picture"].(string)
	for _, v := range oau.RawData["roles"].([]interface{}) {
		if r, ok := v.(string); ok {
			user.Roles = append(user.Roles, r)
		}
	}
	return user, nil
}
