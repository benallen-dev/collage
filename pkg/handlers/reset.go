package handlers

import (
	"net/http"

	"github.com/benallen-dev/collage/pkg/data"
)

func Reset(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value("userData").(*data.SharedData)
	userData.Reset()
}
