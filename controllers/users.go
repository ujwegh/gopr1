package controllers

import (
	"fmt"
	"gopr/models"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "There was an error parsing the form.", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<p>Email: %s</p>", r.FormValue("email"))
	fmt.Fprintf(w, "<p>Name: %s</p>", r.FormValue("name"))
	fmt.Fprintf(w, "<p>Password: %s</p>", r.FormValue("password"))
}
