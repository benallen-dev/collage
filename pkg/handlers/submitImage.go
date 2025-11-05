package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"

	"github.com/benallen-dev/collage/pkg/data"
	"github.com/benallen-dev/collage/pkg/views"
	"github.com/benallen-dev/collage/pkg/util"
)

func SubmitImage(w http.ResponseWriter, r *http.Request) {

	userData := r.Context().Value("userData").(*data.SharedData)

	// This is just a toy app so I'm not going to be elaborate with the error handling.
	name := r.PostFormValue("name")
	sessionId := r.PostFormValue("sessionId")
	file, fileheader, err := r.FormFile("image")
	if err != nil {
		log.Println("error getting file:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check the file is no bigger than 10MB
	// if fileheader.Size > 10*1024*1024 {
	if fileheader.Size > 10*1024*1024 {
		fmt.Println("Uploaded file is too big")
		fmt.Fprintf(w, "<p class=\"error\">Uploaded file is too big. The maximum size is 10MB and your image was " + util.FormatBytes(fileheader.Size) + "</p>")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the file is an image based on MIME type or file extension
	contentType := fileheader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		fmt.Println("Uploaded file is not an image")
		fmt.Fprintf(w, "<p class=\"error\">Uploaded file is not an image</p>")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUuid := uuid.NewString()
	fileExtension := filepath.Ext(fileheader.Filename)
	imagepath := "images/" + newUuid + fileExtension

	// OK, we're all good, let's see if this sessionId already has a file
	// associated with it. If it does, delete it.
	user, ok := userData.GetUser(sessionId)
	if ok && user.ImageUrl != "" {
		if err := os.Remove(user.ImageUrl); err != nil {
			log.Println("error deleting file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Write file to disk
	imageFile, err := os.Create(imagepath)
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

	// Add the user to the shared data
	newUser := data.User{
		Name: name,
		SessionId: sessionId,
		ImageUrl: imagepath,
	}
	
	userData.UpdateUser(newUser)

	// TODO: make a proper template
	// Return HTML containing an img tag with the image in it.

	templ.Handler(views.UpdatedImage(imagepath, name)).ServeHTTP(w, r)
	// fmt.Fprintf(w, "<img className=\"max-w-md\" src=\"%s\" alt=\"%s\" />", imagepath, name)
}
