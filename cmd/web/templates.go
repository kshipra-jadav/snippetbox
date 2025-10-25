package main

import (
	"github.com/kshipra-jadav/snippetbox/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/kshipra-jadav/snippetbox/internal/models"
)

type templateData struct {
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 03:04 PM")
}

var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

func cacheNewTemplate() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		baseName := filepath.Base(page)

		files := []string{
			"html/pages/base.gohtml",
			"html/pages/footer.gohtml",
			"html/pages/nav.gohtml",
			page,
		}

		tmpl, err := template.New(baseName).Funcs(templateFunctions).ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[baseName] = tmpl
	}

	return cache, nil
}
