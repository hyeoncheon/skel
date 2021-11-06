package models_test

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/skel/models"
)

func Test_Doc_String(t *testing.T) {
	r := require.New(t)

	doc := &models.Doc{Title: "Document"}
	r.Equal("Document", doc.String())
}

func Test_Doc_Parent(t *testing.T) {
	r := require.New(t)

	id, err := uuid.FromString("087d1bf4-a4aa-4dfe-9736-8216d1515ca1")
	r.NoError(err)

	_ = models.DB.Destroy(&models.Doc{ID: id})
	doc := &models.Doc{ParentID: id}
	r.Equal(uuid.Nil, doc.Parent().ID)

	err = models.DB.Create(&models.Doc{ID: id, Title: "Parent"})
	r.NoError(err)
	r.Equal(id, doc.Parent().ID)
	r.Equal("Parent", doc.Parent().Title)
}

func Test_Doc_Docs(t *testing.T) {
	r := require.New(t)

	docs := &models.Docs{
		models.Doc{Title: "Great Document"},
		models.Doc{Title: "Poor Document"},
	}
	r.Contains(docs.String(), `"title":"Great Document"`)
	r.Contains(docs.String(), `"title":"Poor Document"`)
	fromJSON := &models.Docs{}
	json.Unmarshal([]byte(docs.String()), fromJSON)
	r.Equal("Great Document", (*fromJSON)[0].Title)
	r.Equal("Poor Document", (*fromJSON)[1].Title)
}

func Test_Doc_Validate(t *testing.T) {
	r := require.New(t)

	doc := &models.Doc{}
	verrs, err := doc.Validate(nil)
	r.NoError(err)
	r.True(verrs.HasAny())
	r.Contains(verrs.String(), "Type can not be")
	r.Contains(verrs.String(), "Category can not be")
	r.Contains(verrs.String(), "Title can not be")
	r.Contains(verrs.String(), "Content can not be")
	r.Contains(verrs.String(), "Permalink can not be")
	r.Contains(verrs.String(), "Lang can not be")

	doc = &models.Doc{
		Type:       "X",
		Category:   "X",
		Title:      "X",
		Content:    "X",
		Permalink:  "X",
		Lang:       "X",
		IsPublic:   true,
		AccessRank: 9,
	}
	verrs, err = doc.Validate(nil)
	r.NoError(err)
	r.False(verrs.HasAny())
	r.Equal(0, doc.AccessRank)
}
