package controllers

import (
	"gopr/views"
	"net/http"
)

type Static struct {
	Template views.Template
}

func (static Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	static.Template.Execute(w, nil)
}

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tpl.Execute(writer, nil)
	}
}
