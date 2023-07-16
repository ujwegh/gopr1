package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to contact page</h1>")
}
func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ page</h1>
		<ul>
			<li>
			<b>Is there free version?</b>
			No only paid version
			</li>
			<li>
			<b>How much does it cost?</b>
			$500
			</li>
			<li>
			<b>Can i pay by credit card?</b>
			Only by cash
			</li>
		</ul>`)
}

type Router struct {
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func main() {
	//http.HandleFunc("/", pathHandler)
	//http.HandleFunc("/contact", contactHandler)
	var router Router
	fmt.Println("arting server on :3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}
