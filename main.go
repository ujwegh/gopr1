package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/home.gohtml")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/contact.gohtml")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/faq.gohtml")
}

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Fprint(w, fmt.Sprintf("hi %v", id))
}

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/galleries/{id}", galleriesHandler)
	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting server on :3000...")
	http.ListenAndServe(":3000", r)
}
