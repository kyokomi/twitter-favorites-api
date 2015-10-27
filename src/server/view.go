package server

import (
	"net/http"

	"github.com/unrolled/render"
)

var renderer = render.New(render.Options{})

type errorView struct {
	Type    string `json:"error_type"`
	Message string `json:"message"`
}

func assembleError(status int, message string) errorView {
	return errorView{
		Type:    ErrorType(status),
		Message: message,
	}
}

func renderErrorResponse(w http.ResponseWriter, status int, msg string) {
	renderer.JSON(w, status, assembleError(status, msg))
}
