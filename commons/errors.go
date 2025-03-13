package commons

import (
	"net/http"
)

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	InternalServerErrorHandler = func(w http.ResponseWriter) {
		WriteJSONError(w, http.StatusInternalServerError, "An unexpected error occurred")
	}
	NotFoundErrorHandler = func(w http.ResponseWriter, err error) {
		WriteJSONError(w, http.StatusNotFound, err.Error())
	}
)