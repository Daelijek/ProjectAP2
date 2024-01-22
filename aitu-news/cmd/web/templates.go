package main

import (
	"html/template"
	"path/filepath"
	"time"
	"alexedwards.net/snippetbox/pkg/forms" // New import
	"alexedwards.net/snippetbox/pkg/models"
)

// Update the templateData fields, removing the individual FormData and
// FormErrors fields and replacing them with a single Form field.
type templateData struct {
	CurrentYear int
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}
...
