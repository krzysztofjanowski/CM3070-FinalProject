package handlers

import (
	"encoding/json"
	"io"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/config"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/render"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/broker"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

// used to store and retrieve all sort of configuration values for Camera Surveillance System app
var applicationConfiguration config.AppConfig

var sessionManager *scs.SessionManager


// this comes from main but need to recreate it to avoid import loop cycle which go does not allow 
func CrossSiteToken(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	return csrfHandler
}

// testing handlers requries a lot of data from main and routes, 
func setupHandlers() http.Handler {

	//We don't need MQTT broker for these tests
	error := broker.Broker()
	if error != nil {
		log.Fatal("MQTT broker connection had one or more failures:", error)
	}

	// 
	applicationConfiguration.TemplateRootDirectory = "../.."
	applicationConfiguration.NotificationsLog = "../../notifications.txt"
	// opening details of the username, email, phone number etc 
	applicationConfiguration.RegistraionDetailsFile = "../../registraionDetailsData.json"

	// check for existing user's details, discard errors as the user may not have registered yet
	registraionDetailsJsonfile, _ := os.Open(applicationConfiguration.RegistraionDetailsFile)

	// opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(registraionDetailsJsonfile)

	var registraionDetails RegistraionDetails

	json.Unmarshal(byteValue, &registraionDetails)

	if registraionDetails.Username != "" {
		// User details already registered so let's copy them to applicationConfig
		applicationConfiguration.UserAlreadyRegistered = true 
		applicationConfiguration.Username = registraionDetails.Username
		applicationConfiguration.Password = registraionDetails.Password 
		applicationConfiguration.Email = registraionDetails.Email
		applicationConfiguration.PhoneNumber = registraionDetails.PhoneNumber
		applicationConfiguration.SlackKey = registraionDetails.SlackKey
	} 

	// video directory
	applicationConfiguration.VideoDir = "../../web/ready_videos"

	// a new session manager 
	sessionManager = scs.New()
	// make the cookies persistes when browser is closed
	sessionManager.Cookie.Persist = true 
	// configure its liftetime this will make the sesion duration to be one day 
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	// pass session to the global applicationConfiguration config so that it can be resused in other places 
	applicationConfiguration.Session = sessionManager
	
	// creating repo to be able to later pass the template cache and other config to the handlers 
	repo := CreateRepo(&applicationConfiguration)
	PassRepo(repo)

	// app config contains template cache and other important config needer in render 
	render.PassConfig(&applicationConfiguration)

	mux := chi.NewRouter()
	
	// to recover from fatal errors 
	mux.Use(middleware.Recoverer)
	mux.Get("/", Repo.IndexPage)
	mux.Get("/starter", Repo.StarterPage)
	mux.Get("/dashboard", Repo.Dashboard)
	mux.Post("/dashboard-post", Repo.DashboardPost)
	mux.Get("/notifications", Repo.Notifications)
	mux.Get("/privacy-on", Repo.PrivacyOn)
	mux.Get("/privacy-off", Repo.PrivacyOff)
	mux.Get("/login", Repo.Login)
	mux.Post("/login-post", Repo.LoginPost)
	mux.Get("/record-now", Repo.RecordNow)
	mux.Post("/registration-post", Repo.RegistrationPost)

	fileServer := http.FileServer(http.Dir("./web/"))

	mux.Handle("/web/*", http.StripPrefix("/web", fileServer))

	return mux 
}
