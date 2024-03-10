package main

import (
	"github.com/go-chi/chi"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/config"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/handlers"
	"net/http"
	"github.com/go-chi/chi/middleware"
)


func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	
	// to recover from fatal errors 
	mux.Use(middleware.Recoverer)

	mux.Use(CrossSiteToken)

	mux.Get("/", handlers.Repo.IndexPage)
	mux.Get("/starter", handlers.Repo.StarterPage)
	mux.Get("/dashboard", handlers.Repo.Dashboard)
	mux.Post("/dashboard-post", handlers.Repo.DashboardPost)
	mux.Get("/notifications", handlers.Repo.Notifications)
	mux.Get("/about", handlers.Repo.AboutGET) // used for dev and testing 
	mux.Get("/privacy-on", handlers.Repo.PrivacyOn)
	mux.Get("/privacy-off", handlers.Repo.PrivacyOff)
	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login-post", handlers.Repo.LoginPost)
	mux.Get("/record-now", handlers.Repo.RecordNow)
	mux.Post("/registration-post", handlers.Repo.RegistrationPost)

	fileServer := http.FileServer(http.Dir("./web/"))
	mux.Handle("/web/*", http.StripPrefix("/web", fileServer))

	return mux 

}