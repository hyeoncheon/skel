package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/skel/models"
)

// UsersResource is the resource for the User model
type UsersResource struct {
	buffalo.Resource
}

// List gets all Users. - GET /users
func (v UsersResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	users := &models.Users{}
	q := tx.PaginateFromParams(c.Params())
	if err := q.All(users); err != nil {
		return errors.WithStack(err)
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.Auto(c, users))
}

// Show gets the data for one User. - GET /users/{user_id}
func (v UsersResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	cu := getCurrentUser(c)
	id := c.Param("user_id")
	isProfileView := false
	if c.Request().URL.Path == "/profile/" {
		id = cu.ID.String()
		isProfileView = true
	}

	user := &models.User{}
	if err := tx.Eager().Find(user, id); err != nil {
		return c.Error(http.StatusNotFound, oops(c, ESU0PS01, err))
	}

	if isProfileView {
		user.Name = cu.Name
		user.Avatar = cu.Avatar
		user.Roles = cu.Roles
	} // else get it from UART

	// CHKME: can I use r.Auto() with c.Set("template")?
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/show"))
}

// Edit renders a edit form for a User. - GET /users/{user_id}/edit
// CHKME: currently not used.
func (v UsersResource) Edit(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.Auto(c, user))
}

// Update changes a User in the DB. - PUT /users/{user_id}
// CHKME: currently not used.
func (v UsersResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.Auto(c, user))
	}

	c.Flash().Add("success", "User was updated successfully")
	return c.Render(http.StatusOK, r.Auto(c, user))
}

// Destroy deletes a User from the DB. - DELETE /users/{user_id}
func (v UsersResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(user); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "User was destroyed successfully")
	return c.Render(http.StatusOK, r.Auto(c, user))
}
