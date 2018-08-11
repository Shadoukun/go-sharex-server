package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	fileReadError        = httpError{"CANT_READ_FILE", http.StatusInternalServerError}
	fileWriteError       = httpError{"CANT_WRITE_FILE", http.StatusInternalServerError}
	invalidFileError     = httpError{"INVALID_FILE", http.StatusBadRequest}
	fileSizeError        = httpError{"FILE_TOO_BIG", http.StatusBadRequest}
	unsupportedFileError = httpError{"UNSUPPORTED_FILE_TYPE", http.StatusBadRequest}
	fileTypeReadError    = httpError{"CANT_READ_FILE_TYPE", http.StatusInternalServerError}
	invalidPasswordError = httpError{"INVALID_PASSWORD", http.StatusForbidden}
)

type httpError struct {
	message    string
	statusCode int
}

func (e *httpError) Error() string {
	return fmt.Sprintf("HTTP %d - %s", e.statusCode, e.message)
}

func renderError(w http.ResponseWriter, err *httpError) {
	w.WriteHeader(err.statusCode)
	w.Write([]byte(err.message + "\n"))
	log.Printf("HTTP CODE: %d - %s", err.statusCode, err.message)
}
