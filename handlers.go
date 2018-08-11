package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type perm interface{}

var (
	userPermission = perm("user-permission")
	tmpl           = template.Must(template.ParseFiles("./static/layout.html"))
)

// pageData stores information for HTML template for Paste uploads.
type pastePageData struct {
	Title string
	Data  string
}

// response is a JSON http response.
type response struct {
	Status int
	Result result
}

type result struct {
	Name      string
	URL       string
	DeleteKey string // unused atm.
}

// uploadPasteHandler is a upload handler for text.
func uploadPasteHandler(w http.ResponseWriter, r *http.Request) {

	// get permissions
	perms := r.Context().Value(userPermission).(*permission)

	r.Body = http.MaxBytesReader(w, r.Body, perms.MaxUpload<<20)
	fmt.Println(r)
	if err := r.ParseMultipartForm(perms.MaxUpload << 20); err != nil {
		fmt.Println(err)
		renderError(w, &fileSizeError)
		return
	}

	// parse and validate file and post parameters
	paste, _, err := r.FormFile("file")
	if err != nil {
		renderError(w, &fileReadError)
		fmt.Println(err)
		return
	}
	fileBytes, err := ioutil.ReadAll(paste)
	if err != nil {
		renderError(w, &fileReadError)
		fmt.Println(err)
		return
	}

	// creates a buffer and compiles data for template.
	tpl := bytes.Buffer{}
	data := pastePageData{
		Title: "Paste",
		Data:  string(fileBytes),
	}

	// renders the  template.
	if err := tmpl.Execute(&tpl, data); err != nil {
		fmt.Println(err)
		renderError(w, &fileWriteError)
		return
	}

	// Create file in save location.
	newFilename := randomToken(6)
	newPath := filepath.Join(uploadPastePath, newFilename)
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, &fileWriteError)
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(tpl.Bytes()); err != nil || newFile.Close() != nil {
		renderError(w, &fileWriteError)
		return
	}

	// send the response
	w.Write(createResponse("paste", newFilename))
}

// uploadFileHandler is a upload handler for files.
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	// get permissions
	perms := r.Context().Value(userPermission).(*permission)

	// validate file size
	r.Body = http.MaxBytesReader(w, r.Body, perms.MaxUpload<<20)
	if err := r.ParseMultipartForm(perms.MaxUpload << 20); err != nil {
		renderError(w, &fileSizeError)
		fmt.Println(perms.MaxUpload)
		return
	}

	// parse and validate file and post parameters
	file, _, err := r.FormFile("file")
	if err != nil {
		renderError(w, &invalidFileError)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(w, &invalidFileError)
		return
	}

	// check file type.
	fileType, err := validateFileType(&fileBytes, perms)
	if err != nil {
		renderError(w, err.(*httpError))
		return
	}

	// Create file in save location.
	newFilename := randomToken(6) + fileType
	newPath := filepath.Join(uploadFilePath, newFilename)
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, &fileWriteError)
		return
	}
	defer newFile.Close() // idempotent, okay to call twice

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderError(w, &fileWriteError)
		return
	}

	// send the response
	w.Write(createResponse("file", newFilename))
}

// permissionMiddleware validates the password provided in a request.
// If the password is valid, it adds the correct permissionGroup to the
// request's context and passes it to the upload handler.
func permissionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx context.Context
		var pass = r.FormValue("pass")

		switch {
		case pass == config.Admin.Password:
			fmt.Println("Admin Upload")
			ctx = context.WithValue(ctx, userPermission, &config.Admin)
			next.ServeHTTP(w, r.WithContext(ctx))

		case pass == config.Default.Password:
			fmt.Println("Default Upload")
			ctx = context.WithValue(ctx, userPermission, &config.Default)
			next.ServeHTTP(w, r.WithContext(ctx))

		default:
			renderError(w, &invalidPasswordError)
			return
		}
	})
}

// detailsHandler responds with simple file details.
func detailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte(vars["file"]))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/pages/404.html")

}
