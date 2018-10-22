package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Doc is a structure for storing document data
type Doc struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Type        string    `json:"type" db:"type"`
	Category    string    `json:"category" db:"category"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Permalink   string    `json:"permalink" db:"permalink"`
	Lang        string    `json:"lang" db:"lang"`
	AccessRank  int       `json:"access_rank" db:"access_rank"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	IsPublished bool      `json:"is_published" db:"is_published"`
	AuthorID    uuid.UUID `json:"author_id" db:"author_id"`
	ParentID    uuid.UUID `json:"parent_id" db:"parent_id"`
}

// String represents doc as a string
func (d Doc) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

//*** docs ---

// Docs is an array of docs
type Docs []Doc

// String represents docs as a string
func (d Docs) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

//*** validation functions ---

// Validate gets run every time you call a "pop.Validate*" method.
func (d *Doc) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Type, Name: "Type"},
		&validators.StringIsPresent{Field: d.Category, Name: "Category"},
		&validators.StringIsPresent{Field: d.Title, Name: "Title"},
		&validators.StringIsPresent{Field: d.Content, Name: "Content"},
		&validators.StringIsPresent{Field: d.Permalink, Name: "Permalink"},
		&validators.StringIsPresent{Field: d.Lang, Name: "Lang"},
		&validators.IntIsPresent{Field: d.AccessRank, Name: "AccessRank"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (d *Doc) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (d *Doc) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
