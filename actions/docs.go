package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/skel/models"
)

// DocsResource is the resource for the Doc model
type DocsResource struct {
	buffalo.Resource
}

// List gets all Docs. - GET /docs
func (v DocsResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	docs := &models.Docs{}

	if !getCurrentUser(c).IsAdmin() {
		if err := tx.
			Eager("Author").
			Eager("Children.Author").
			Eager("Children.Children.Author").
			Where("parent_id = ?", uuid.Nil).All(docs); err != nil {
			return errors.WithStack(err)
		}
		c.Set("pagination", tx.PaginateFromParams(c.Params()).Paginator)
		return c.Render(http.StatusOK, r.Auto(c, docs))
	}

	q := tx.PaginateFromParams(c.Params())
	if err := q.Eager().All(docs); err != nil {
		return errors.WithStack(err)
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.Auto(c, docs))
}

// Show gets the data for one Doc. - GET /docs/{doc_id}
func (v DocsResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	doc := &models.Doc{}
	if err := tx.
		Eager("Author").
		Eager("Children.Author").
		Eager("Children.Children.Author").
		Where("lang = ?", currentLanguage(c)).
		Where("permalink = ?", c.Param("doc_id")).
		First(doc); err != nil {
		if err := tx.
			Eager("Author").
			Eager("Children.Author").
			Eager("Children.Children.Author").
			Where("lang = ?", defaultLanguage).
			Where("permalink = ?", c.Param("doc_id")).
			First(doc); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	}

	return c.Render(http.StatusOK, r.Auto(c, doc))
}

// ShowByLang gets the data for one Doc. - GET /docs/{lang}/{doc_permalink}
func (v DocsResource) ShowByLang(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	doc := &models.Doc{}
	if err := tx.
		Eager("Author").
		Eager("Children.Author").
		Eager("Children.Children.Author").
		Where("lang = ?", c.Param("lang")).
		Where("permalink = ?", c.Param("permalink")).
		First(doc); err != nil {
		c.Flash().Add("success", "redirected")
		return c.Redirect(http.StatusTemporaryRedirect, "docPath()",
			map[string]interface{}{
				"doc_id": c.Param("permalink"),
			})
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
		return errors.WithStack(errors.New("no transaction found"))
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
		return errors.WithStack(errors.New("no transaction found"))
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
		return errors.WithStack(errors.New("no transaction found"))
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

// Publish updates publishing status of the doc. - GET /docs/{doc_id}/publish
func (v DocsResource) Publish(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	doc := &models.Doc{}
	if err := tx.Find(doc, c.Param("doc_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	user := getCurrentUser(c)
	if !user.IsAdmin() && doc.AuthorID != user.ID {
		return c.Render(http.StatusForbidden, r.HTML("violation"))
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
