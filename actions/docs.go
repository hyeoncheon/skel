package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/skel/models"
)

// DocsResource is the resource for the Doc model
type DocsResource struct {
	buffalo.Resource
}

func init() {
	addPermission("GET newDocsPath", "doctor")
	addPermission("POST docsPath", "doctor")
	addPermission("GET editDocPath", "doctor")
	addPermission("PUT docPath", "doctor")
	addPermission("DELETE docPath", "doctor")
	addPermission("PUT docPublishPath", "doctor")
}

func addPermissionStatement(q *pop.Query, user *models.User) {
	if user.IsAdmin() {
		return
	} else if user.HasRole("doctor") {
		q = q.Where("(is_published = ? OR author_id = ?)", true, user.ID)
	} else {
		q = q.Where("is_published = ?", true)
		if user.IsGuest() {
			q = q.Where("is_public = ?", true)
		}
	}
}

// List gets all Docs. - GET /docs
func (v DocsResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	docs := &models.Docs{}
	q := tx.PaginateFromParams(c.Params()).Eager("Author")
	user := getCurrentUser(c)
	if !user.IsAdmin() {
		c.Set("view", "tree")
		q = q.
			Eager("Children.Author").
			Eager("Children.Children.Author").
			Where("parent_id = ?", uuid.Nil)
		addPermissionStatement(q, user)
	}

	if err := q.Order("permalink").All(docs); err != nil {
		return errors.WithStack(err) // TODO prettify
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.Auto(c, docs))
}

// Show gets the data for one Doc. - GET /docs/{doc_id}
func (v DocsResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}
	c.Logger().Debugf("doc.show:%v %v", c.Param("doc_id"), currentLanguage(c))

	doc := &models.Doc{}
	permalink := c.Param("doc_id")
	q := tx.Eager("Author").
		Eager("Children.Author").
		Eager("Children.Children.Author").
		Where("permalink = ?", permalink)
	addPermissionStatement(q, getCurrentUser(c))

	cq := pop.Q(q.Connection)
	q.Clone(cq)
	if err := cq.Where("lang = ?", currentLanguage(c)).First(doc); err != nil {
		c.Logger().Warnf("doc:%v has no %v version", permalink, currentLanguage(c))
		if err := q.Where("lang = ?", defaultLanguage).First(doc); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	}

	return c.Render(http.StatusOK, r.Auto(c, doc))
}

// ShowByLang gets the data for one Doc. - GET /docs/{lang}/{doc_permalink}
func (v DocsResource) ShowByLang(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}
	c.Logger().Debugf("doc.show:%v @%v", c.Param("permalink"), c.Param("lang"))

	doc := &models.Doc{}
	permalink := c.Param("permalink")
	q := tx.Eager("Author").
		Eager("Children.Author").
		Eager("Children.Children.Author").
		Where("permalink = ?", permalink)
	addPermissionStatement(q, getCurrentUser(c))

	if err := q.Where("lang = ?", c.Param("lang")).First(doc); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "docPath()",
			render.Data{"doc_id": permalink})
	}

	return c.Render(http.StatusOK, r.Auto(c, doc))
}

// New renders the form for creating a new Doc. - GET /docs/new
func (v DocsResource) New(c buffalo.Context) error {
	doc := &models.Doc{}
	doc.ParentID, _ = uuid.FromString(c.Param("parent"))
	return c.Render(http.StatusOK, r.Auto(c, doc))
}

// Create adds a Doc to the DB. - POST /docs
func (v DocsResource) Create(c buffalo.Context) error {
	doc := &models.Doc{}
	if err := c.Bind(doc); err != nil {
		return errors.WithStack(err)
	}

	doc.AuthorID = getCurrentUser(c).ID
	if doc.Lang == "" {
		doc.Lang = currentLanguage(c)
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	verrs, err := tx.ValidateAndCreate(doc)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.Auto(c, doc))
	}

	c.Flash().Add("success", "Doc was created successfully")
	return c.Redirect(http.StatusSeeOther, "docLangPath()",
		map[string]interface{}{"lang": doc.Lang, "permalink": doc.Permalink})
}

// Edit renders a edit form for a Doc. - GET /docs/{doc_id}/edit
func (v DocsResource) Edit(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	doc := &models.Doc{}
	if err := tx.Eager("Author").
		Find(doc, c.Param("doc_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.Auto(c, doc))
}

// Update changes a Doc in the DB. - PUT /docs/{doc_id}
func (v DocsResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	doc := &models.Doc{}
	if err := tx.Find(doc, c.Param("doc_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(doc); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(doc)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.Auto(c, doc))
	}

	c.Flash().Add("success", "Doc was updated successfully")
	return c.Redirect(http.StatusSeeOther, "docLangPath()",
		map[string]interface{}{"lang": doc.Lang, "permalink": doc.Permalink})
}

// Destroy deletes a Doc from the DB. - DELETE /docs/{doc_id}
func (v DocsResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	doc := &models.Doc{}
	if err := tx.Find(doc, c.Param("doc_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(doc); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "Doc was destroyed successfully")
	return c.Render(http.StatusOK, r.Auto(c, doc))
}

// Publish updates publishing status of the doc. - PUT /docs/{doc_id}/publish
func (v DocsResource) Publish(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return oops(c, ESTX0001, nil)
	}

	doc := &models.Doc{}
	if err := tx.Find(doc, c.Param("doc_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	doc.IsPublished = !doc.IsPublished

	verrs, err := tx.ValidateAndUpdate(doc)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.Auto(c, doc))
	}

	c.Flash().Add("success", "Doc was updated successfully")
	return c.Redirect(http.StatusSeeOther, "docLangPath()",
		map[string]interface{}{"lang": doc.Lang, "permalink": doc.Permalink})
}
