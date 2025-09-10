package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/kshipra-jadav/snippetbox/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 03:04 PM")
}

var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

func cacheNewTemplate() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("../../ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		baseName := filepath.Base(page)

		files := []string{
			"../../ui/html/pages/base.html",
			"../../ui/html/pages/footer.html",
			page,
		}

		tmpl, err := template.New(baseName).Funcs(templateFunctions).ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[baseName] = tmpl
	}

	return cache, nil
}
