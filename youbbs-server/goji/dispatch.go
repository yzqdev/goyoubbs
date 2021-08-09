package goji

import (
	"net/http"

	"goyoubbs/goji/internal"
)

type dispatch struct{}

func (d dispatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h := ctx.Value(internal.Handler)
	if h == nil {
		http.NotFound(w, r)
	} else {
		h.(http.Handler).ServeHTTP(w, r)
	}
}
