package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/vusile/misa-saa-ngapi/handler"
	"github.com/vusile/misa-saa-ngapi/model"
	"github.com/vusile/misa-saa-ngapi/repository"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	if getDotEnvValue("IS_PRODUCTION") == "false" {
		router.Use(middleware.Logger,
			csrf.Protect([]byte("32-byte-long-auth-key"),
				csrf.Path("/"),
				csrf.Secure(false)))
	} else {
		router.Use(middleware.Logger,
			csrf.Protect([]byte("32-byte-long-auth-key"),
				csrf.Path("/")))
	}

	homeHandler := &handler.HomeHandler{
		Client:   a.gorm,
		ESClient: a.esClient,
	}

	adminHandler := &handler.AdminHandler{
		Client:   a.gorm,
		ESClient: a.esClient,
	}

	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("/go/src/app/assets"))))
	router.Get("/", homeHandler.Home)
	router.Get("/admin", adminHandler.Home)
	router.Post("/search", homeHandler.Search)
	router.Route("/users", a.loadUserRoutes)
	router.Route("/majimbo", a.loadJimboRoutes)
	router.Route("/admin/majimbo", a.loadAdminJimboRoutes)
	router.Route("/countries", a.loadCountryRoutes)
	router.Route("/churches", a.loadChurchRoutes)
	router.Route("/parokia", a.loadParokiaRoutes)
	router.Route("/languages", a.loadLanguagesRoutes)
	router.Route("/timings", a.loadTimingRoutes)

	a.router = router
}

func (a *App) loadUserRoutes(router chi.Router) {
	userHandler := &handler.User{Repo: &repository.UserRepo{
		Client: a.gorm,
	}}

	router.Group(func(router chi.Router) {
		router.Use(a.loggedInMiddleware)
		router.Get("/login", userHandler.LoginForm)
		router.Post("/login", userHandler.Login)
		router.Get("/register", userHandler.RegistrationForm)
		router.Post("/", userHandler.Create)
		router.Get("/confirm-account/{id}", userHandler.CodeForm)
		router.Post("/confirm", userHandler.ConfirmAccount)
	})

	router.Group(func(router chi.Router) {
		router.Use(a.authMiddleware)
		router.Get("/logout", userHandler.Logout)
		router.Get("/", userHandler.List)
		router.Get("/{id}", userHandler.GetByID)
		router.Put("/{id}", userHandler.UpdateByID)
		router.Delete("/{id}", userHandler.DeleteByID)
	})
}

func (a *App) loadJimboRoutes(router chi.Router) {
	jimboHandler := &handler.Jimbo{Repo: &repository.JimboRepo{
		Client: a.gorm,
	}}
	router.Get("/", jimboHandler.All)
	router.Get("/{slug}/{id}", jimboHandler.Detail)

}

func (a *App) loadAdminJimboRoutes(router chi.Router) {
	jimboHandler := &handler.Jimbo{Repo: &repository.JimboRepo{
		Client: a.gorm,
	}}

	router.Use(a.authMiddleware)
	router.Post("/", jimboHandler.Create)
	router.Get("/", jimboHandler.List)
	router.Get("/{id}", jimboHandler.GetByID)
	router.Put("/{id}", jimboHandler.UpdateByID)
	router.Delete("/{id}", jimboHandler.DeleteByID)

}

func (a *App) loadCountryRoutes(router chi.Router) {
	countryHandler := &handler.Country{Repo: &repository.CountryRepo{
		Client: a.gorm,
	}}

	router.Use(a.authMiddleware)
	router.Post("/", countryHandler.Create)
	router.Get("/", countryHandler.List)
	router.Get("/{id}", countryHandler.GetByID)
	router.Put("/{id}", countryHandler.UpdateByID)
	router.Delete("/{id}", countryHandler.DeleteByID)

}

func (a *App) loadChurchRoutes(router chi.Router) {
	churchHandler := &handler.Church{Repo: &repository.ChurchRepo{
		Client: a.gorm,
	}}

	router.Post("/", churchHandler.Create)
	router.Get("/", churchHandler.List)
	router.Get("/{id}", churchHandler.GetByID)
	router.Put("/{id}", churchHandler.UpdateByID)
	router.Delete("/{id}", churchHandler.DeleteByID)

}

func (a *App) loadParokiaRoutes(router chi.Router) {
	parokiaHandler := &handler.Parokia{Repo: &repository.ParokiaRepo{
		Client:   a.gorm,
		ESClient: a.esClient,
	}}

	router.Get("/{slug}/{id}", parokiaHandler.Detail)

	router.Group(func(router chi.Router) {
		router.Use(a.authMiddleware)
		router.Post("/", parokiaHandler.Create)
		router.Get("/", parokiaHandler.List)
		router.Get("/{id}", parokiaHandler.GetByID)
		router.Put("/{id}", parokiaHandler.UpdateByID)
		router.Delete("/{id}", parokiaHandler.DeleteByID)
	})
}

func (a *App) loadLanguagesRoutes(router chi.Router) {
	languageHandler := &handler.Language{Repo: &repository.LanguageRepo{
		Client: a.gorm,
	}}

	router.Post("/", languageHandler.Create)
	router.Get("/", languageHandler.List)
	router.Get("/{id}", languageHandler.GetByID)
	router.Put("/{id}", languageHandler.UpdateByID)
	router.Delete("/{id}", languageHandler.DeleteByID)

}

func (a *App) loadTimingRoutes(router chi.Router) {
	timingHandler := &handler.Timing{Repo: &repository.TimingRepo{
		Client: a.gorm,
	}}

	router.Use(a.authMiddleware)
	router.Get("/timingform/{id}", timingHandler.List)
	router.Post("/", timingHandler.Create)
	router.Get("/", timingHandler.List)
	router.Get("/{id}", timingHandler.GetByID)
	router.Put("/{id}", timingHandler.UpdateByID)
	router.Get("/delete/{id}/{parokiaId}", timingHandler.DeleteByID)

}

func (a *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st, err := r.Cookie("session_token")
		user := model.User{}

		if err != nil || st.Value == "" {
			http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
		} else {
			csrf, err := r.Cookie("csrf_token")

			if csrf.Value != "" && err == nil {
				a.gorm.Where("session_token = ?", st.Value).First(&user)

				if csrf.Value != user.CsrfToken {
					http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
				}
			} else {
				http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) loggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st, err := r.Cookie("session_token")
		user := model.User{}

		if err == nil && st.Value != "" {
			csrf, err := r.Cookie("csrf_token")

			if csrf.Value != "" && err == nil {
				a.gorm.Where("session_token = ?", st.Value).First(&user)

				if csrf.Value == user.CsrfToken {
					http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
