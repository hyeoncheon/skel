package actions

import (
	"net/http"
	"testing"

	"github.com/hyeoncheon/skel/models"
)

func (as *ActionSuite) Test_UsersResource_List() {
	as.NoError(as.setupUsers())

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/users").Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
		res = as.HTML(res.Location()).Get()
		as.Contains(res.Body.String(), "Login required")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/users").Get()
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		res := as.HTML("/users").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), users["admin"].Email)
		as.Contains(res.Body.String(), users["doctor"].Email)
		as.Contains(res.Body.String(), users["user"].Email)
	}))
}

func (as *ActionSuite) Test_UsersResource_Show() {
	as.NoError(as.setupUsers())
	id := as.userID("user")

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/users/%v", id).Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/users/%v", id).Get()
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		res := as.HTML("/users/%v", id).Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.NotContains(res.Body.String(), users["admin"].Email)
		as.NotContains(res.Body.String(), users["doctor"].Email)
		as.Contains(res.Body.String(), users["user"].Email)
	}))
}

func (as *ActionSuite) Test_UsersResource_Profile() {
	as.NoError(as.setupUsers())

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/profile").Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/profile").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), users["user"].Email)
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		res := as.HTML("/profile").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), users["admin"].Email)
	}))
}

func (as *ActionSuite) Test_UsersResource_Edit() {
	as.T().Log("UsersResource.Edit is not used currently")
	as.True(true)
}

func (as *ActionSuite) Test_UsersResource_Update() {
	as.T().Log("UsersResource.Update is not used currently")
	as.True(true)
}

func (as *ActionSuite) Test_UsersResource_Destroy() {
	as.NoError(as.setupUsers())
	id := as.userID("doctor")

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/users/%v", id).Delete()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/users/%v", id).Delete()
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		res := as.HTML("/users/%v", id).Delete()
		as.Equal(http.StatusFound, res.Code, "invalid status")
		as.Equal("/users/", res.Location(), "invalid redirection")
		as.NotContains(res.Body.String(), users["doctor"].Email)
		cnt, err := models.DB.Count(&models.User{})
		as.NoError(err)
		as.Equal(3, cnt)
	}))
}
