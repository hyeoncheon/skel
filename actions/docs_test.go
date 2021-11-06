package actions

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/skel/models"
)

func (as *ActionSuite) Test_DocsResource_List() {
	as.NoError(as.setupUsers())
	as.NoError(as.setupDocs())

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) { // none
		res := as.HTML("/docs").Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asGuest", func(t *testing.T) { // only public published
		as.NoError(as.asUser("guest"))

		res := as.HTML("/docs").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.NotContains(res.Body.String(), ">Published<")
		as.NotContains(res.Body.String(), ">Public Draft<")
		as.Contains(res.Body.String(), ">Public Published<")
	}))
	as.True(t.Run("asUser", func(t *testing.T) { // published
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Published<")
		as.NotContains(res.Body.String(), ">Public Draft<")
		as.Contains(res.Body.String(), ">Public Published<")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) { // all
		as.NoError(as.asUser("doctor"))

		res := as.HTML("/docs").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Published<")
		as.Contains(res.Body.String(), ">Public Draft<")
		as.Contains(res.Body.String(), ">Public Published<")
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) { // all
		as.NoError(as.asUser("admin"))

		res := as.HTML("/docs").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Published<")
		as.Contains(res.Body.String(), ">Public Draft<")
		as.Contains(res.Body.String(), ">Public Published<")
	}))
}

func (as *ActionSuite) Test_DocsResource_Show() {
	as.NoError(as.setupUsers())
	as.NoError(as.setupDocs())

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) { // none
		res := as.HTML("/docs/%v", "published").Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")

		res = as.HTML("/docs/%v", "public-draft").Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")

		res = as.HTML("/docs/%v", "public-published").Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asGuest", func(t *testing.T) { // only public-published
		as.NoError(as.asUser("guest"))

		res := as.HTML("/docs/%v", "published").Get()
		as.Equal(http.StatusNotFound, res.Code, "invalid status")

		res = as.HTML("/docs/%v", "public-draft").Get()
		as.Equal(http.StatusNotFound, res.Code, "invalid status")

		res = as.HTML("/docs/%v", "public-published").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Public Published<")
	}))
	as.True(t.Run("asUser", func(t *testing.T) { // only published
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs/%v", "published").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Published<")

		res = as.HTML("/docs/%v", "public-draft").Get()
		as.Equal(http.StatusNotFound, res.Code, "invalid status")

		res = as.HTML("/docs/%v", "public-published").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Public Published<")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) { // published and mine
		as.NoError(as.asUser("doctor"))

		res := as.HTML("/docs/%v", "published").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Published<")

		res = as.HTML("/docs/%v", "public-draft").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Public Draft<")

		res = as.HTML("/docs/%v", "public-published").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Published<")
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) { // all
		as.NoError(as.asUser("admin"))

		res := as.HTML("/docs/%v", "published").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Published<")

		res = as.HTML("/docs/%v", "public-draft").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Public Draft<")

		res = as.HTML("/docs/%v", "public-published").Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), ">Public Published<")
	}))
}

func (as *ActionSuite) Test_DocsResource_ShowByLang() {
	as.NoError(as.setupUsers())
	as.NoError(as.setupDocs())
	as.NoError(as.asUser("user"))

	res := as.HTML("/docs/ko/%v", "language").Get()
	as.Equal(http.StatusOK, res.Code, "invalid status")
	as.Contains(res.Body.String(), ">Language<")
	as.Contains(res.Body.String(), "한글문서")

	res = as.HTML("/docs/en/%v", "language").Get()
	as.Equal(http.StatusOK, res.Code, "invalid status")
	as.Contains(res.Body.String(), ">Language<")
	as.Contains(res.Body.String(), "English Document")

	res = as.HTML("/docs/ja/%v", "language").Get()
	as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
	as.Equal("/docs/language/", res.Location())
}

func (as *ActionSuite) Test_DocsResource_New() {
	as.NoError(as.setupUsers())

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) { // none
		res := as.HTML("/docs/new").Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) { // none
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs/new").Get()
		as.Equal(http.StatusForbidden, res.Code, "invalide status")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) { // none
		as.NoError(as.asUser("doctor"))

		res := as.HTML("/docs/new").Get()
		as.Equal(http.StatusOK, res.Code, "invalide status")
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) { // none
		as.NoError(as.asUser("admin"))

		res := as.HTML("/docs/new").Get()
		as.Equal(http.StatusOK, res.Code, "invalide status")
	}))
}

func (as *ActionSuite) Test_DocsResource_Create() {
	as.NoError(as.setupUsers())
	invalidDoc := map[string]interface{}{"IsPublic": false}
	validDoc := map[string]interface{}{
		"Title":    "title of doc",
		"Type":     "type of doc",
		"Category": "category of doc",
		"Content":  "content of doc",
	}

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/docs").Post(validDoc)
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs").Post(validDoc)
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) {
		as.NoError(as.asUser("doctor"))

		res := as.HTML("/docs").Post(invalidDoc)
		as.Equal(http.StatusUnprocessableEntity, res.Code, "invalid status")
		as.Contains(res.Body.String(), "invalid-feedback")

		res = as.HTML("/docs").Post(validDoc)
		as.Equal(http.StatusSeeOther, res.Code, "invalid status")
		as.Equal("/docs/en/title-of-doc/", res.Location())
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		validDoc["Lang"] = "ko"
		res := as.HTML("/docs").Post(validDoc)
		as.Equal(http.StatusSeeOther, res.Code, "invalid status")
		as.Equal("/docs/ko/title-of-doc/", res.Location())
	}))
}

func (as *ActionSuite) Test_DocsResource_Edit() {
	as.NoError(as.setupUsers())
	as.NoError(as.setupDocs())

	doc := as.getDoc("top1")
	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/docs/%v/edit", doc.ID).Get()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs/%v/edit", doc.ID).Get()
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) {
		as.NoError(as.asUser("doctor"))

		res := as.HTML("/docs/%v/edit", doc.ID).Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), docs["top1"].Content)
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		res := as.HTML("/docs/%v/edit", doc.ID).Get()
		as.Equal(http.StatusOK, res.Code, "invalid status")
		as.Contains(res.Body.String(), docs["top1"].Content)
	}))
}

func (as *ActionSuite) Test_DocsResource_Update() {
	as.NoError(as.setupUsers())
	as.NoError(as.setupDocs())

	data := map[string]interface{}{"IsPublic": true}
	data2 := map[string]interface{}{"IsPublic": false}

	doc := as.getDoc("top2")
	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/docs/%v", doc.ID).Put(data)
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs/%v", doc.ID).Put(data)
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) {
		as.NoError(as.asUser("doctor"))

		as.Equal(false, doc.IsPublic)
		res := as.HTML("/docs/%v", doc.ID).Put(data)
		as.Equal(http.StatusSeeOther, res.Code, "invalid status")
		as.Equal(true, as.getDoc("top2").IsPublic)
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		as.Equal(true, as.getDoc("top2").IsPublic)
		res := as.HTML("/docs/%v", doc.ID).Put(data2)
		as.Equal(http.StatusSeeOther, res.Code, "invalid status")
		as.Equal(false, as.getDoc("top2").IsPublic)
	}))
}

func (as *ActionSuite) Test_DocsResource_Destroy() {
	as.NoError(as.setupUsers())
	as.NoError(as.setupDocs())

	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/docs/%v", as.getDoc("top2").ID).Delete()
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs/%v", as.getDoc("top2").ID).Delete()
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) {
		as.NoError(as.asUser("doctor"))

		res := as.HTML("/docs/%v", as.getDoc("top2").ID).Delete()
		as.Equal(http.StatusFound, res.Code, "invalid status")
		as.Equal(uuid.Nil, as.getDoc("top2").ID)
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		res := as.HTML("/docs/%v", as.getDoc("top1").ID).Delete()
		as.Equal(http.StatusFound, res.Code, "invalid status")
		as.Equal(uuid.Nil, as.getDoc("top1").ID)
	}))
}

func (as *ActionSuite) Test_DocsResource_Publish() {
	as.NoError(as.setupUsers())
	as.NoError(as.setupDocs())

	as.Equal(false, as.getDoc("top2").IsPublished)
	t := as.T()
	as.True(t.Run("asNobody", func(t *testing.T) {
		res := as.HTML("/docs/%v/publish", as.getDoc("top2").ID).Put(nil)
		as.Equal(http.StatusTemporaryRedirect, res.Code, "invalid status")
		as.Equal("/", res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asUser", func(t *testing.T) {
		as.NoError(as.asUser("user"))

		res := as.HTML("/docs/%v/publish", as.getDoc("top2").ID).Put(nil)
		as.Equal(http.StatusForbidden, res.Code, "invalid status")
		as.Contains(res.Body.String(), "VIOLATION")
	}))
	as.True(t.Run("asDoctor", func(t *testing.T) {
		as.NoError(as.asUser("doctor"))

		res := as.HTML("/docs/%v/publish", as.getDoc("top2").ID).Put(nil)
		as.Equal(http.StatusSeeOther, res.Code, "invalid status")
		doc := as.getDoc("top2")
		as.Equal(true, doc.IsPublished)
		permalink := fmt.Sprintf("/docs/%v/%v/", doc.Lang, doc.Permalink)
		as.Equal(permalink, res.Location(), "invalid redirection")
	}))
	as.True(t.Run("asAdmin", func(t *testing.T) {
		as.NoError(as.asUser("admin"))

		doc := as.getDoc("top2")
		as.Equal(true, doc.IsPublished)
		res := as.HTML("/docs/%v/publish", doc.ID).Put(nil)
		as.Equal(http.StatusSeeOther, res.Code, "invalid status")

		doc = as.getDoc("top2")
		as.Equal(false, doc.IsPublished)
		permalink := fmt.Sprintf("/docs/%v/%v/", doc.Lang, doc.Permalink)
		as.Equal(permalink, res.Location(), "invalid redirection")
	}))
}

//** helpers ---

func (as *ActionSuite) getDoc(k string) *models.Doc {
	doc := &models.Doc{}
	_ = models.DB.Where("title = ?", docs[k].Title).First(doc)
	return doc
}

func (as *ActionSuite) setupDocs() error {
	doctor := &models.User{}
	as.NoError(models.DB.Where("email = ?", users["doctor"].Email).First(doctor))
	for _, doc := range docs {
		doc.AuthorID = doctor.ID
		children := doc.Children
		verrs, err := models.DB.ValidateAndCreate(&doc)
		as.NoError(errors.WithStack(err))
		as.False(verrs.HasAny(), "validation failed: %v", verrs)
		for _, child := range children {
			child.AuthorID = doctor.ID
			verrs, err := models.DB.ValidateAndCreate(&child)
			as.NoError(errors.WithStack(err))
			as.False(verrs.HasAny(), "validation failed: %v", verrs)
		}
	}
	count, err := models.DB.Count(&models.Doc{})
	as.NoError(err)
	as.Equal(9, count)
	return nil
}

var docs = map[string]models.Doc{
	"top1": models.Doc{
		Title:       "Top level document #1",
		Lang:        "en",
		Type:        "Manual",
		Category:    "Test",
		Content:     "content for document #1",
		IsPublished: true,
		Children: models.Docs{
			models.Doc{
				Title:    "Second level document #1.1",
				Lang:     "en",
				Type:     "Manual",
				Category: "Test",
				Content:  "content for document #1.1",
			},
			models.Doc{
				Title:    "Second level document #1.2",
				Lang:     "en",
				Type:     "Manual",
				Category: "Test",
				Content:  "content for document #1.2",
			},
		},
	},
	"top2": models.Doc{
		Title:    "Top level document #2",
		Lang:     "en",
		Type:     "Manual",
		Category: "Test",
		Content:  "content for document #2",
	},
	"public-draft": models.Doc{
		Title:       "Public Draft",
		Lang:        "en",
		Type:        "Manual",
		Category:    "Test",
		Content:     "content for document #2",
		IsPublic:    true,
		IsPublished: false,
	},
	"public-published": models.Doc{
		Title:       "Public Published",
		Lang:        "en",
		Type:        "Manual",
		Category:    "Test",
		Content:     "content for document #2",
		IsPublic:    true,
		IsPublished: true,
	},
	"published": models.Doc{
		Title:       "Published",
		Lang:        "en",
		Type:        "Manual",
		Category:    "Test",
		Content:     "content for document #2",
		IsPublic:    false,
		IsPublished: true,
	},
	"en-doc": models.Doc{
		Title:       "Language",
		Lang:        "en",
		Type:        "Manual",
		Category:    "Test",
		Content:     "English Document",
		IsPublic:    false,
		IsPublished: true,
	},
	"ko-doc": models.Doc{
		Title:       "Language",
		Lang:        "ko",
		Type:        "Manual",
		Category:    "Test",
		Content:     "한글문서",
		IsPublic:    false,
		IsPublished: true,
	},
}
