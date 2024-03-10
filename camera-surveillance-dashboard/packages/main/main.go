package main

import (
	"encoding/json"
	"fmt"
	"io"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/broker"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/config"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/handlers"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

// network port number for the dashboard ui
const portNumber string = ":8080"

// used to store and retrieve all sort of configuration values for Camera Surveillance System app
var applicationConfiguration config.AppConfig

var sessionManager *scs.SessionManager

func initiateBackend() error {

	//MQTT broker
	error := broker.Broker()
	if error != nil {
		log.Fatal("MQTT broker connection had one or more failures:", error)
	}
	
	applicationConfiguration.TemplateRootDirectory = "."

	// opening details of the username, email, phone number etc 
	applicationConfiguration.RegistraionDetailsFile = "registraionDetailsData.json"

	// check for existing user's details, discard errors as the user may not have registered yet
	registraionDetailsJsonfile, _ := os.Open(applicationConfiguration.RegistraionDetailsFile)

	// opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(registraionDetailsJsonfile)

	var registraionDetails handlers.RegistraionDetails

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
	applicationConfiguration.VideoDir = "web/ready_videos"
	applicationConfiguration.NotificationsLog = "notifications.txt"
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
	repo := handlers.CreateRepo(&applicationConfiguration)
	handlers.PassRepo(repo)

	// app config contains template cache and other important config needer in render 
	render.PassConfig(&applicationConfiguration)

	return nil
}

func initateWebServer() error {

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&applicationConfiguration),
	}

	httpServerError := server.ListenAndServe()
	if httpServerError != nil {
		log.Fatal(httpServerError)
		return httpServerError
	}

	return nil 
}

func main() {

	fmt.Println("Starting Camera Surveillance backend")

	backendError  := initiateBackend()
	if backendError != nil {
		log.Fatal("Camera surveillance backend failed to initate") 
	}

	fmt.Printf("Starting Camera Surveillance System dashboard at port%s \n", portNumber)
	
	webServerError := initateWebServer()
	if webServerError != nil {
		log.Fatal("Camera surveillance webserver failed to initate") 
	}

}
