package controllers

import "net/http"

type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) // What arguments does this take?
}
