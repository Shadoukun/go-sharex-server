package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const uploadFilePath = "./data/files"
const uploadPastePath = "./data/pastes"
const staticFilePath = "./static"

func main() {
	log.SetOutput(os.Stdout)
	config.getConf()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// Static files
	staticfs := http.FileServer(http.Dir(staticFilePath))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticfs))

	// Paste
	pastefs := http.FileServer(http.Dir(uploadPastePath))

	psr := r.PathPrefix("/paste").Subrouter()
	r.Handle("/paste/{file}", http.StripPrefix("/paste", pastefs))
	psr.HandleFunc("/upload", permissionMiddleware(uploadPasteHandler))
	psr.HandleFunc("/{file}/details", detailsHandler)

	// File
	filefs := http.FileServer(http.Dir(uploadFilePath))

	fsr := r.PathPrefix("/file").Subrouter()
	r.Handle("/file/{file}", http.StripPrefix("/file", filefs))
	fsr.HandleFunc("/{file}/details", detailsHandler)
	fsr.HandleFunc("/upload", permissionMiddleware(uploadFileHandler))

	log.Print("Server started on " + config.Host + ":" + config.Port)
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, r))
}
