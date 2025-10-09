package httputils

import (
	"fmt"
	"net/http"
)

// error messages
const (
	ErrorBadRequest    = "BAD_REQUEST"
	ErrorUnathorized   = "UNATHORIZED"
	ErrorForbidden     = "FORBIDDEN"
	ErrorNotFound      = "NOT_FOUND"
	ErrorConflict      = "CONFLICT"
	ErrorInternalError = "INTERNAL_ERROR"
)

// BadRequestJSON template for bad request status response
func BadRequestJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{ "error": "BAD_REQUEST" }`)
}

// ForbiddenJSON template for forbidden status response
func ForbiddenJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, `{ "error": "FORBIDDEN" }`)
}

// NotFoundJSON template for not found status response
func NotFoundJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{ "error": "NOT_FOUND" }`)
}

// ConflictJSON template for conflict status response
func ConflictJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusConflict)
	fmt.Fprintf(w, `{ "error": "CONFLICT" }`)
}

// InternalErrorJSON template for internal error status response
func InternalErrorJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{ "error": "INTERNAL_ERROR" }`)
}
