package views

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:generate templ generate

// FsBootstrap is an embed.FS that contains the bootstrap files
//
//go:embed bootstrap-5.1.3-dist/*
var FsBootstrap embed.FS

// BootStrapHandler is a function that returns a http.HandlerFunc
func BootStrapHandler(subPathStripPrefix string) http.HandlerFunc {
	bootstrapFS, err := fs.Sub(FsBootstrap, "bootstrap-5.1.3-dist")
	if err != nil {
		panic(fmt.Errorf("failed getting the sub tree for the site files: %w", err))
	}
	handler := http.FileServer(http.FS(bootstrapFS))
	static := http.StripPrefix(subPathStripPrefix, handler)
	return func(w http.ResponseWriter, r *http.Request) {
		static.ServeHTTP(w, r)
	}
}
