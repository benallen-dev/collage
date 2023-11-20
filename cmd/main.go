package main

import (
	"context"
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

	r.Get("/presenter", func (w http.ResponseWriter, r *http.Request) {
		// This is probably not the best way to do this but meh whatever IT WORKS BABY

		// Get the shared data from the request context
		userData := r.Context().Value("userData").(*data.SharedData)
		// Get the users from the shared data
		users := userData.GetUsers()
		// Get the presenter view
		presenterView := views.Presenter(users)
		// Render the presenter view
		templ.Handler(presenterView).ServeHTTP(w, r)	
	})

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
			users := r.Context().Value("userData").(*data.SharedData).GetUsers()
			fmt.Fprintf(w, "%v", users)
		})
	})

	log.Println("Listening on :1323...")
	http.ListenAndServe(":1323", r)
}
