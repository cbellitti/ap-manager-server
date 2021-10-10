package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// A (simple) example of our application-wide configuration.
type Env struct {
	DB *sql.DB
	//SessionStorage *storage.Session
	//DBRepo     *mysqldb.MysqlDBRepo
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	*Env
	H      func(e *Env, w http.ResponseWriter, r *http.Request) error
	Method string
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != h.Method {
		http.Error(w, GetErrorMessageForCode(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func HandlerParseResponse(w http.ResponseWriter, r []byte, e error) error {
	if e != nil {
		//http.Error(w, e.Error(), http.StatusInternalServerError)
		return e
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(r)
	return nil
}

var ERRORS map[string]string = make(map[string]string)

func GetStatusErrorForCode(code int) StatusError {
	return StatusError{code, errors.New(ERRORS[strconv.Itoa(code)])}
}

func GetErrorMessageForCode(code int) string {
	return ERRORS[strconv.Itoa(code)]
}

func GetErrorForCode(code int) error {
	return errors.New(ERRORS[strconv.Itoa(code)])
}

func SetErrors() {
	ERRORS[strconv.Itoa(http.StatusInternalServerError)] = "Internal Server Error"
	ERRORS[strconv.Itoa(http.StatusUnauthorized)] = "Unauthorized"
	ERRORS[strconv.Itoa(http.StatusBadRequest)] = "Bad Request"
	ERRORS[strconv.Itoa(http.StatusMethodNotAllowed)] = "Method Not Allowed"
}
