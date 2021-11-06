package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/hyeoncheon/skel/models"
)

// getCurrentUser is wrapper function for getCurrentUser family.
func getCurrentUser(c buffalo.Context) *models.User {
	return getCurrentUserSimple(c)
}

// getCurrentUserSimple is simplified version of getCurrentUser family.
// It simeply build the user instance with session values without
// database access.
func getCurrentUserSimple(c buffalo.Context) *models.User {
	user := &models.User{}
	if id, ok := c.Session().Get("user_id").(uuid.UUID); ok {
		user.ID = id
		user.Email = c.Session().Get("user_mail").(string)
		user.Name = c.Session().Get("user_name").(string)
		user.Avatar = c.Session().Get("user_icon").(string)
		user.Roles = c.Session().Get("user_roles").([]string)
	}
	return user
}

// getCurrentUserDatabase finds user from database backend then
// returns the user instance after filling it with session information.
func getCurrentUserDatabase(c buffalo.Context) *models.User {
	user := &models.User{}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		c.Logger().Error("system error: no transaction found")
		return user
	}

	if err := tx.Find(user, c.Session().Get("user_id")); err == nil {
		user.Name = c.Session().Get("user_name").(string)
		user.Avatar = c.Session().Get("user_icon").(string)
		user.Roles = c.Session().Get("user_roles").([]string)
	}
	return user
}

// makeSession registers user information in session.
// this function should be called by authenication/authorization process.
func makeSession(c buffalo.Context, user *models.User) error {
	sess := c.Session()
	sess.Set("user_id", user.ID)
	sess.Set("user_mail", user.Email)
	sess.Set("user_name", user.Name)
	sess.Set("user_icon", user.Avatar)
	sess.Set("user_roles", user.Roles)
	return sess.Save()
}

// destroySession clears all session values.
// this function must be called by logout handler.
func destroySession(c buffalo.Context) error {
	c.Session().Clear()
	return nil
}
