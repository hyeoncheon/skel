package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// User is a structure for storing service users
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"-"`
	Avatar    string    `json:"avatar" db:"-"`
	Roles     []string  `json:"roles" db:"-"`
	Docs      Docs      `has_many:"docs" fk_id:"author_id" order_by:"permalink"`
}

// String represents user as a string (currently returns email address)
func (u User) String() string {
	if u.Name != "" {
		return fmt.Sprintf("%s (%s)", u.Name, u.ID.String()[0:6])
	}
	return fmt.Sprintf("%s (%s)", u.Email, u.ID.String()[0:6])
}

// IsValid returns true if the user instance has valid values
// but it does not check database entries for validation.
func (u *User) IsValid() bool {
	if u.ID == uuid.Nil {
		return false
	}
	return true
}

// IsAdmin returns true if the user instance has role `admin` by calling
// `User.HasRole("admin")`.
func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}

// HasRole returns true if the user instance contains given `role` value
// in its `User.Roles` array.
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsGuest returns true if the user has only "guest" role.
func (u *User) IsGuest() bool {
	if len(u.Roles) == 1 && u.Roles[0] == "guest" {
		return true
	}
	return false
}

//*** users ---

// Users is an array of users
type Users []User

// String represents users as a string
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

//*** validation functions ---

// Validate gets run every time you call a "pop.Validate*"
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	if u.Email == "" {
		return nil, errors.New("Invalid user: Email is not provided")
	}
	if u.ID == uuid.Nil {
		return nil, errors.New("Invalid user: Invalid user ID")
	}
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.ID, Name: "ID"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
