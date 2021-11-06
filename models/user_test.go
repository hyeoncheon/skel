package models_test

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/skel/models"
)

func Test_User_String(t *testing.T) {
	r := require.New(t)

	user := &models.User{Name: "Yonghwan"}
	r.Equal("Yonghwan (000000)", user.String())

	user = &models.User{Email: "sio4@example.com"}
	r.Equal("sio4@example.com (000000)", user.String())
}

func Test_User_IsValid(t *testing.T) {
	r := require.New(t)

	user := &models.User{}
	r.False(user.IsValid())
	var err error
	user.ID, err = uuid.FromString("087d1bf4-a4aa-4dfe-9736-8216d1515ca1")
	r.NoError(err)
	r.True(user.IsValid())
}

func Test_User_IsAdmin(t *testing.T) {
	r := require.New(t)

	user := &models.User{}
	user.Roles = []string{"user"}
	r.False(user.IsAdmin())

	user.Roles = append(user.Roles, "admin")
	r.True(user.IsAdmin())
}

func Test_User_IsGuest(t *testing.T) {
	r := require.New(t)

	user := &models.User{}
	user.Roles = []string{"guest"}
	r.True(user.IsGuest())

	user.Roles = append(user.Roles, "user")
	r.False(user.IsGuest())
}

func Test_User_Users(t *testing.T) {
	r := require.New(t)

	users := &models.Users{
		models.User{Name: "Yonghwan"},
		models.User{Name: "Stark"},
	}
	r.Contains(users.String(), `"name":"Yonghwan"`)
	r.Contains(users.String(), `"name":"Stark"`)
	fromJSON := &models.Users{}
	json.Unmarshal([]byte(users.String()), fromJSON)
	r.Equal("Yonghwan", (*fromJSON)[0].Name)
	r.Equal("Stark", (*fromJSON)[1].Name)
}

func Test_User_Validate(t *testing.T) {
	r := require.New(t)

	user := &models.User{}
	_, err := user.Validate(nil)
	r.Error(err)
	r.Contains(err.Error(), "Email is not")

	user.Email = "sio4@example.com"
	_, err = user.Validate(nil)
	r.Error(err)
	r.Contains(err.Error(), "Invalid user ID")

	user.ID, err = uuid.FromString("087d1bf4-a4aa-4dfe-9736-8216d1515ca1")
	r.NoError(err)
	_, err = user.Validate(nil)
	r.NoError(err)
}
