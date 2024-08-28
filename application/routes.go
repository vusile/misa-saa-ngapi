package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vusile/misa-saa-ngapi/handler"
	"github.com/vusile/misa-saa-ngapi/repository"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/users", a.loadUserRoutes)
	router.Route("/invoices", a.loadInvoiceRoutes)
	router.Route("/promises", a.loadPromisesRoutes)
	router.Route("/majimbo", a.loadJimboRoutes)
	router.Route("/countries", a.loadCountryRoutes)
	router.Route("/churches", a.loadChurchRoutes)
	router.Route("/parokia", a.loadParokiaRoutes)
	router.Route("/languages", a.loadLanguagesRoutes)

	a.router = router
}

func (a *App) loadInvoiceRoutes(router chi.Router) {
	invoiceHandler := &handler.Invoice{}

	router.Post("/", invoiceHandler.Create)
	router.Get("/", invoiceHandler.List)
	router.Get("/{id}", invoiceHandler.GetByID)
	router.Put("/{id}", invoiceHandler.UpdateByID)
	router.Delete("/{id}", invoiceHandler.DeleteByID)
}

func (a *App) loadUserRoutes(router chi.Router) {
	userHandler := &handler.User{}

	router.Post("/", userHandler.Create)
	router.Get("/", userHandler.List)
	router.Get("/{id}", userHandler.GetByID)
	router.Put("/{id}", userHandler.UpdateByID)
	router.Delete("/{id}", userHandler.DeleteByID)
}

func (a *App) loadPromisesRoutes(router chi.Router) {
	promiseHandler := &handler.Promise{}

	router.Post("/", promiseHandler.Create)
	router.Get("/", promiseHandler.List)
	router.Get("/{id}", promiseHandler.GetByID)
	router.Put("/{id}", promiseHandler.UpdateByID)
	router.Delete("/{id}", promiseHandler.DeleteByID)

}

func (a *App) loadJimboRoutes(router chi.Router) {
	jimboHandler := &handler.Jimbo{Repo: &repository.JimboRepo{
		Client: a.gorm,
	}}

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
		Client: a.gorm,
	}}

	router.Post("/", parokiaHandler.Create)
	router.Get("/", parokiaHandler.List)
	router.Get("/{id}", parokiaHandler.GetByID)
	router.Put("/{id}", parokiaHandler.UpdateByID)
	router.Delete("/{id}", parokiaHandler.DeleteByID)

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
