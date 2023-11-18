package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/benallen-dev/collage/pkg/views"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/", templ.Handler(views.Index()).ServeHTTP)
	r.Get("/presenter", templ.Handler(views.Presenter()).ServeHTTP)
	r.Get("/participant", templ.Handler(views.Participant()).ServeHTTP)

	// Serve static images
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "images"))
	FileServer(r, "/images", filesDir)

	r.Route("/api", func(r chi.Router) {
		r.Post("/submit", handleSubmit)
	})

	log.Println("Listening on :1323...")
	http.ListenAndServe(":1323", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// 
// Blatantly stolen from https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
// TODO: Move this to a separate package/file
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

// TODO: Move this to a separate package/file
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	// This is just a toy app so I'm not going to be elaborate with the error handling.
	name := r.PostFormValue("name")
	file, fileheader, err := r.FormFile("image")
	if err != nil {
		log.Println("error getting file:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Println("name:", name)
	log.Println("file:", fileheader.Filename)

	// Check if the file is an image based on MIME type or file extension
	contentType := fileheader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		fmt.Println("Uploaded file is not an image")
		fmt.Fprintf(w, "<p class=\"error\">Uploaded file is not an image</p>")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Write file to disk
	imageFile, err := os.Create("images/" + fileheader.Filename)
	if err != nil {
		log.Println("error creating file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer imageFile.Close()

	if _, err := io.Copy(imageFile, file); err != nil {
		log.Println("error copying file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return HTML containing an img tag with the image in it.
	fmt.Fprintf(w, "<img src=\"/images/%s\" alt=\"%s\" />", fileheader.Filename, name)
}
