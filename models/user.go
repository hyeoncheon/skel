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

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"-"`
	Avatar    string    `json:"avatar" db:"-"`
	Roles     []string  `json:"roles" db:"-"`
}

// String represents user as string (currently returns email address)
func (u User) String() string {
	return fmt.Sprintf("%s (%s)", u.Email, u.ID.String()[0:6])
}

func (u *User) IsValid() bool {
	if u.ID == uuid.Nil {
		return false
	}
	return true
}

//*** users ---

// Users is an array of users
type Users []User

// String represents users as string
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
