package controllers

import "net/http"

type Template interface {
	Execute(w http.ResponseWriter, data interface{}) // What arguments does this take?
}
