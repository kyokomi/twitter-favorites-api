package server

import (
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func NewAppEngineHttpClient(ctx appengine.Context) *http.Client {
	return urlfetch.Client(ctx)
}
