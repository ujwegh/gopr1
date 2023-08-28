package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"gopr/controllers"
	"gopr/migrations"
	"gopr/models"
	"gopr/templates"
	"gopr/views"
	"net/http"
	"time"
)

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Fprint(w, fmt.Sprintf("hi %v", id))
}

func TimerMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		fmt.Println("Request time:", time.Since(start))
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	// Setup a database connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// Setup our model services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}
	// Setup our controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	tpl := views.MustParse(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.MustParse(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.MustParse(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	usersC.Templates.New = views.MustParse(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.MustParse(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.Me = views.MustParse(views.ParseFS(templates.FS, "me.gohtml", "tailwind.gohtml"))

	r.Get("/signup", TimerMiddleware(usersC.New))
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)

	r.Get("/galleries/{id}", galleriesHandler)
	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "Page not found", http.StatusNotFound)
	})

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
	)
	handler := csrfMw(umw.SetUser(r))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", handler)
}
