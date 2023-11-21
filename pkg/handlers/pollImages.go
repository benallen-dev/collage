package handlers

import (
	"net/http"
	"sort"

	"github.com/a-h/templ"

	"github.com/benallen-dev/collage/pkg/data"
	"github.com/benallen-dev/collage/pkg/views"
)

func PollImages(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value("userData").(*data.SharedData)
	users := userData.GetUsers()

	sort.Slice(users, func(i, j int) bool {
		return users[i].SessionId < users[j].SessionId
	})

	// Render the images
	templ.Handler(views.Images(users)).ServeHTTP(w, r)	
}
