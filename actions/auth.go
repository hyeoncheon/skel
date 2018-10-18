package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
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
		c.Logger().Warnf("could not set user: %v", err)
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
	if !ok || len(roles) < 0 {
		return errors.New("Invalid user: You have no access permission")
	}
	return nil
}

// setUser creates or updates user in database, and returns user.
func setUser(c buffalo.Context, oau *goth.User) (*models.User, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		c.Logger().Error("system error: no transaction found")
		return &models.User{}, errors.New("Internal error: No transaction")
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
		return &models.User{}, err
	}
	if verrs.HasAny() {
		return &models.User{}, errors.New("Internal error: Validation failed")
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

func makeSession(c buffalo.Context, user *models.User) error {
	sess := c.Session()
	sess.Set("user_id", user.ID)
	sess.Set("user_mail", user.Email)
	sess.Set("user_name", user.Name)
	sess.Set("user_icon", user.Avatar)
	sess.Set("user_roles", user.Roles)
	return sess.Save()
}

func destroySession(c buffalo.Context) error {
	c.Session().Clear()
	return nil
}
