package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
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
	Author      User      `json:"author" belongs_to:"user"`
	ParentID    uuid.UUID `json:"parent_id" db:"parent_id"`
	Children    Docs      `json:"children" has_many:"docs" fk_id:"parent_id" order_by:"permalink"`
}

// String represents doc as a string
func (d Doc) String() string {
	return d.Title
}

// Parent returns parent document of the doc.
func (d *Doc) Parent() *Doc {
	doc := &Doc{}
	if d.ParentID != uuid.Nil {
		if err := DB.Find(doc, d.ParentID); err != nil {
			// TODO: log here
		}
	}
	return doc
}

func (d *Doc) permalink() string {
	return url.PathEscape(flect.Dasherize(d.Title))
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
	if d.Permalink == "" {
		d.Permalink = d.permalink()
	}
	if d.IsPublic {
		d.AccessRank = 0
	}
	fmt.Printf("----------------- %v\n", d.permalink())
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Type, Name: "Type"},
		&validators.StringIsPresent{Field: d.Category, Name: "Category"},
		&validators.StringIsPresent{Field: d.Title, Name: "Title"},
		&validators.StringIsPresent{Field: d.Content, Name: "Content"},
		&validators.StringIsPresent{Field: d.Permalink, Name: "Permalink"},
		&validators.StringIsPresent{Field: d.Lang, Name: "Lang"},
	), nil
}
