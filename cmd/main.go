package main

import (
	"context"
	"log"
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

// Middleware to add a shared variable to request context
func curryMiddleware(userData *data.SharedData) func(next http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add your shared variable here (e.g., sharedVar := "someValue")
			// Create a context with the shared variable
			ctx := context.WithValue(r.Context(), "userData" , userData)
	
			// Pass the context with the shared variable to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func main() {
	
	userData := data.NewSharedData()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(curryMiddleware(userData)) // now we should have userData on the r.Context() object

	r.Get("/", templ.Handler(views.Index()).ServeHTTP)
	r.Get("/participant", templ.Handler(views.Participant(uuid.NewString())).ServeHTTP)
	r.Get("/presenter", templ.Handler(views.Presenter()).ServeHTTP)

	// Serve static images
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "images"))
	staticDir := http.Dir(filepath.Join(workDir, "static"))

	handlers.FileServer(r, "/", staticDir) // Used for css
	handlers.FileServer(r, "/images", filesDir) // Used for images

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/poll", handlers.PollImages)
		
		r.Post("/submit", handlers.SubmitImage)

		// r.Post("/reset", handlers.Reset)
		// r.Post("/delete/{sessionId}", handlers.Delete)
	})

	log.Println("Listening on :1323...")
	http.ListenAndServe(":1323", r)
}
