package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Using "all:" ensures hidden files/folders (if any) are included
//
//go:embed all:dist
var embeddedFiles embed.FS

func GetFileSystem() http.FileSystem {
	f, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		panic("cannot find subtree of filesystem: " + err.Error())
	}
	return http.FS(f)
}

func ServeFrontend(r chi.Router) {

	frontendFilesystem := GetFileSystem()
	fileServer := http.FileServer(frontendFilesystem)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/api") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "API route not found"}`))
			return
		}

		// Clean the path to check the internal filesystem
		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		// Search for the static asset (js, css, favicon, etc.)
		file, err := frontendFilesystem.Open(filePath)
		if err == nil {
			file.Close()
			fileServer.ServeHTTP(w, r) // File exists, serve it directly
			return
		}

		// File does not exist -> This is an SPA route. Serve index.html instead.
		indexFile, err := frontendFilesystem.Open("index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		stat, _ := indexFile.Stat()
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile)

	})

}
