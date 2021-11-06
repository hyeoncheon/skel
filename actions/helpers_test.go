package actions

import (
	"github.com/gofrs/uuid"

	"github.com/hyeoncheon/skel/models"
)

//* basic test helpers
func (as *ActionSuite) asUser(k string) error { // same as makeSession()
	user := &models.User{}
	err := models.DB.Where("email = ?", users[k].Email).First(user)
	as.NoError(err)

	as.Session.Set("user_id", user.ID)
	as.Session.Set("user_mail", user.Email)
	as.Session.Set("user_name", users[k].Name)
	as.Session.Set("user_icon", users[k].Avatar)
	as.Session.Set("user_roles", users[k].Roles)

	return nil
}

func (as *ActionSuite) userID(k string) uuid.UUID {
	user := &models.User{}
	err := models.DB.Where("email = ?", users[k].Email).First(user)
	as.NoError(err)
	return user.ID
}

func (as *ActionSuite) setupUsers() error {
	for _, user := range users {
		err := models.DB.Create(&user)
		as.NoError(err)
	}
	return nil
}

var users = map[string]models.User{
	"admin": models.User{
		Name:   "Admininstrator",
		Email:  "admin@example.com",
		Avatar: "http://example.com/icon.png",
		Roles:  []string{"admin"},
	},
	"doctor": models.User{
		Name:   "Documentor",
		Email:  "doctor@example.com",
		Avatar: "http://example.com/icon.png",
		Roles:  []string{"doctor", "user"},
	},
	"user": models.User{
		Name:   "Normal User",
		Email:  "user@example.com",
		Avatar: "http://example.com/icon.png",
		Roles:  []string{"user"},
	},
	"guest": models.User{
		Name:   "Guest User",
		Email:  "guest@example.com",
		Avatar: "http://example.com/icon.png",
		Roles:  []string{"guest"},
	},
}
