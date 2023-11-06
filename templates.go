package main

import (
	"embed"
	"html/template"
)

var (
	//go:embed templates
	tmplfs embed.FS
	tmpl   = template.Must(template.ParseFS(tmplfs, "templates/*.gohtml"))
)
