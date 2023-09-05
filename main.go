package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"gopr/controllers"
	"gopr/migrations"
	"gopr/models"
	"gopr/templates"
	"gopr/views"
	"net/http"
	"os"
	"strconv"
)

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Fprint(w, fmt.Sprintf("hi %v", id))
}

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.DefaultPostgresConfig()
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}
	//cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	//cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	// Setup our config
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	// Setup the database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Migrate the database
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// Setup our model services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{
		DB: db,
	}
	// Setup our controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)
	handler := csrfMw(umw.SetUser(r))

	tpl := views.MustParse(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.MustParse(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.MustParse(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	usersC.Templates.New = views.MustParse(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.MustParse(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.Me = views.MustParse(views.ParseFS(templates.FS, "me.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.MustParse(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	usersC.Templates.CheckYourEmail = views.MustParse(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ResetPassword = views.MustParse(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))
	galleriesC.Templates.New = views.MustParse(views.ParseFS(templates.FS, "galleries/new.gohtml", "tailwind.gohtml"))

	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	r.Get("/galleries/{id}", galleriesHandler)
	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "Page not found", http.StatusNotFound)
	})
	r.Get("/galleries/new", galleriesC.New)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	// Start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, handler)
	if err != nil {
		panic(err)
	}
}
