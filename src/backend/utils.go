// This file provides some utility types and functions used by the rest of the
// code of the codelab.
// You don't need to modify anything in this file.

package flicq

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"appengine"
	"appengine/user"
)

// appError is an error with a HTTP response code.
type appError struct {
	error
	Code int
}

// appErrorf creates a new appError given a reponse code and a message.
func appErrorf(code int, format string, args ...interface{}) *appError {
	return &appError{fmt.Errorf(format, args...), code}
}

// appHandler handles HTTP requests and manages returned errors.
type appHandler func(w io.Writer, r *http.Request) error

// appHandler implements http.Handler.
func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	buf := &bytes.Buffer{}
	err := h(buf, r)
	if err == nil {
		io.Copy(w, buf)
		return
	}
	code := http.StatusInternalServerError
	logf := c.Errorf
	if err, ok := err.(*appError); ok {
		code = err.Code
		logf = c.Infof
	}

	w.WriteHeader(code)
	logf(err.Error())
	fmt.Fprint(w, err)
}

// authReq checks that a user is logged in before executing the appHandler.
type authReq appHandler

// authReq implements http.Handler.
func (h authReq) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if user.Current(c) == nil {
		http.Error(w, "login required", http.StatusForbidden)
		return
	}
	appHandler(h).ServeHTTP(w, r)
}
