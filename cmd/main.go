package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/benallen-dev/collage/pkg/handlers"
	"github.com/benallen-dev/collage/pkg/views"
	"github.com/benallen-dev/collage/pkg/data"
)

func init() {

	// Create the images directory if it doesn't exist
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		os.Mkdir("images", 0755)
	}

	// If it does, empty it
	files, err := os.ReadDir("images")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if err := os.Remove("images/" + file.Name()); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	
	userData := data.NewSharedData()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/", templ.Handler(views.Index()).ServeHTTP)
	r.Get("/presenter", templ.Handler(views.Presenter(userData)).ServeHTTP)
	r.Get("/participant", templ.Handler(views.Participant(uuid.NewString())).ServeHTTP)

	// Serve static images
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "images"))
	handlers.FileServer(r, "/images", filesDir)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/submit", func (w http.ResponseWriter, r *http.Request) {
			handlers.SubmitImage(w, r, userData)
		})

		r.Get("/users", func (w http.ResponseWriter, r *http.Request) {
			users := userData.GetUsers()
			fmt.Fprintf(w, "%v", users)
		})
	})

	log.Println("Listening on :1323...")
	http.ListenAndServe(":1323", r)
}
