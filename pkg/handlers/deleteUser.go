package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/a-h/templ"

	"github.com/benallen-dev/collage/pkg/data"
	"github.com/benallen-dev/collage/pkg/views"

	"github.com/go-chi/chi/v5"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE USER BUTTON PRESSED")
	sessionId := chi.URLParam(r, "sessionId")

	userData := r.Context().Value("userData").(*data.SharedData)
	userData.DeleteUser(sessionId)

	users := userData.GetUsers()
	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})

	templ.Handler(views.Images(users)).ServeHTTP(w, r)	
}
