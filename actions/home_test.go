package actions

import (
	"net/http"
)

func (as *ActionSuite) Test_HomeHandler_Without_Login() {
	res := as.HTML("/").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "Are you an Hyeoncheon member")
}

func (as *ActionSuite) Test_HomeHandler_With_Login() {
	as.NoError(as.setupUsers())
	as.NoError(as.asUser("user"))

	res := as.HTML("/").Get()
	as.Equal(http.StatusTemporaryRedirect, res.Code)
	as.Equal("/profile/", res.Location())
}

func (as *ActionSuite) Test_LogoutHandler() {
	as.NoError(as.setupUsers())
	as.NoError(as.asUser("user"))

	as.NotNil(as.Session.Get("user_id"))
	res := as.HTML("/logout").Get()
	as.Equal(http.StatusTemporaryRedirect, res.Code)
	as.Equal("/", res.Location())
	as.Nil(as.Session.Get("user_id"))
}
