package main

import (
	"net/http"
	"log"

	"github.com/benallen-dev/collage/pkg/util"
	"github.com/benallen-dev/collage/pkg/views"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})

	http.HandleFunc("/name", func(w http.ResponseWriter, r *http.Request) {
		views.Hello(util.GetRandomName()).Render(r.Context(), w)
	})

	log.Println("Listening on :1323...")
	err := http.ListenAndServe(":1323", nil)
	if err != nil {
		log.Fatal(err);
	}
}
