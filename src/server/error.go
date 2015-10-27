package server

import "net/http"

var errorTypes = map[int]string{
	http.StatusBadRequest:          "bad_request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
	http.StatusNotFound:            "not_found",
	http.StatusConflict:            "conflict",
	http.StatusInternalServerError: "internal_server_error",
}

func ErrorType(status int) string {
	return errorTypes[status]
}
