package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
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
	q := tx.PaginateFromParams(c.Params())
	if err := q.All(docs); err != nil {
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
	if err := tx.Find(doc, c.Param("doc_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.Auto(c, doc))
}

// New renders the form for creating a new Doc. - GET /docs/new
func (v DocsResource) New(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.Auto(c, &models.Doc{}))
}

// Create adds a Doc to the DB. - POST /docs
func (v DocsResource) Create(c buffalo.Context) error {
	doc := &models.Doc{}
	if err := c.Bind(doc); err != nil {
		return errors.WithStack(err)
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
	return c.Render(http.StatusCreated, r.Auto(c, doc))
}

// Edit renders a edit form for a Doc. - GET /docs/{doc_id}/edit
func (v DocsResource) Edit(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	doc := &models.Doc{}
	if err := tx.Find(doc, c.Param("doc_id")); err != nil {
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
	return c.Render(http.StatusOK, r.Auto(c, doc))
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
