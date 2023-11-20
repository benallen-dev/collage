package handlers

import (
	"fmt"
	"net/http"
)

func Reset(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RESET BUTTON PRESSED")
}
