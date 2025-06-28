package assets

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed static/*
var staticEmbed embed.FS

//go:embed templates/*
var templateEmbed embed.FS

var (
	staticFS   fs.FS
	templateFS fs.FS
)

func init() {
	var err error

	staticFS, err = fs.Sub(staticEmbed, "static")
	if err != nil {
		log.Fatalf("failed to sub static embed: %v", err)
	}

	templateFS, err = fs.Sub(templateEmbed, "templates")
	if err != nil {
		log.Fatalf("failed to sub templates embed: %v", err)
	}
}

func StaticFS() fs.FS {
	return staticFS
}

func TemplateFS() fs.FS {
	return templateFS
}
