package handlers

import (
	"encoding/json"
	"fmt"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/broker"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/config"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/helpers"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/models"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/render"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/videos"
	"log"
	"net/http"
	"os"
	"strings"
	"sort"
)

// app config bucket 
type Repository struct {
	applicationConfiguration *config.AppConfig
}

//  handlers repo 
var Repo *Repository

// creates a pointer to the repository to pass data from other packages 
func CreateRepo(a *config.AppConfig) *Repository {
	return &Repository{
		applicationConfiguration: a,
	}
}

// allows passing of data to the repository for the handlers
func PassRepo(r *Repository) {
	Repo = r
}

func (repo *Repository) AboutGET(w http.ResponseWriter, req *http.Request) {

	render.RenderWebPage(w, req, "about.page.tmpl", &models.WebData{
	})

}


// to display history of notification messages sent
func (repo *Repository) Notifications(w http.ResponseWriter, req *http.Request) {

	var Notifications []string

	// Read the content of the file 'notifications.txt'
	content, err := os.ReadFile(Repo.applicationConfiguration.NotificationsLog)
	if err != nil {
		// fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Convert the content to string and split it into lines
	lines := strings.Split(string(content), "\n")

	// Iterate over the lines and add them to the Notifications slice
	Notifications = append(Notifications, lines...)

	render.RenderWebPage(w, req, "notifications.page.tmpl", &models.WebData{Notifications: Notifications})
}

func (repo *Repository) IndexPage(w http.ResponseWriter, req *http.Request) {

	// logical 'OR' is needed when registraiton form is filled first time
	if Repo.applicationConfiguration.UserAlreadyRegistered ||  Repo.applicationConfiguration.Username != "" {
		render.RenderWebPage(w, req, "login.page.tmpl", &models.WebData{})
	}else {
		render.RenderWebPage(w, req, "starter.page.tmpl", &models.WebData{})
	}
}

func (repo *Repository) StarterPage(w http.ResponseWriter, req *http.Request) {

	render.RenderWebPage(w, req, "starter.page.tmpl", &models.WebData{})
}

// main dashboard with all recorded videos 
func (repo *Repository) Dashboard(w http.ResponseWriter, req *http.Request) {
	// fmt.Println("applicationConfiguration.VideoDir:",Repo.applicationConfiguration.VideoDir)

	videos,err  := videos.ListVideos(Repo.applicationConfiguration.VideoDir)
	if err != nil {
		log.Fatal("fatal error im porting videos:", err, videos)
		return
	}

	sort.Slice(videos, func(i, j int) bool {
        return videos[i] < videos[j]
    })
	
	render.RenderWebPage(w, req, "dashboard.page.tmpl", 
							&models.WebData{
											LatestVideo: videos[len(videos)-1],
											MovementSensor1 : broker.WebData.MovementSensor1 , 
											MovementSensor2 : false, LightSensor: 
											broker.WebData.LightSensor, 
											VideosSlice: videos,
											})

}


// main dashboard with all recorded videos 
func (repo *Repository) DashboardPost(w http.ResponseWriter, req *http.Request) {

	startDateOnForm := req.Form.Get("ffromDate")
	endDateOnForm   := req.Form.Get("ftoDate")
	
	filteredVideos, err := helpers.FindMatchingVideos(startDateOnForm , endDateOnForm, Repo.applicationConfiguration.VideoDir)
	
	fmt.Println("len:", len(filteredVideos))

	sort.Slice(filteredVideos, func(i, j int) bool {
        return filteredVideos[i] < filteredVideos[j]
    })

	if err != nil {
		log.Fatal("fatal error finding matching videos:", err)
		return
	}

	render.RenderWebPage(w, req, "dashboard.page.tmpl", 
						 &models.WebData{LatestVideo: filteredVideos[len(filteredVideos)-1], 
										 MovementSensor1 : broker.WebData.MovementSensor1 , 
										 MovementSensor2 : false, 
										 LightSensor: broker.WebData.LightSensor, 
										 VideosSlice: filteredVideos})
}


// Used to enable privacy
func (repo *Repository) PrivacyOn(w http.ResponseWriter, req *http.Request) {

	broker.PublishToTopic(broker.Client,"camera/video", "DISABLE ALL RECORDING")
	
	render.RenderWebPage(w, req, "dashboard.page.tmpl", &models.WebData{})
}


// Used to disable privacy
func (repo *Repository) PrivacyOff(w http.ResponseWriter, req *http.Request) {

	broker.PublishToTopic(broker.Client,"camera/video", "ENABLE ALL RECORDING")

	render.RenderWebPage(w, req, "dashboard.page.tmpl", &models.WebData{})
}

// Used to record now
func (repo *Repository) RecordNow(w http.ResponseWriter, req *http.Request) {

	// pulish to MQTT to start a new recording 
	broker.PublishToTopic(broker.Client,"pir_sensor/motion_reading", "movement_detected-manual")

	render.RenderWebPage(w, req, "dashboard.page.tmpl", &models.WebData{})
}


// Login page 
func (repo *Repository) Login(w http.ResponseWriter, req *http.Request) {

	render.RenderWebPage(w, req, "login.page.tmpl", &models.WebData{})
}

// Login page 
func (repo *Repository) LoginPost(w http.ResponseWriter, req *http.Request) {

	username := req.Form.Get("fusername")
	password := req.Form.Get("fpassword")

	if username == Repo.applicationConfiguration.Username && password == Repo.applicationConfiguration.Password {
		// User authenticated
		videos,err  := videos.ListVideos(Repo.applicationConfiguration.VideoDir)
		if err != nil {
			log.Fatal("fatal error importing videos:", err, videos)
			return
		}
		
		render.RenderWebPage(w, req, "dashboard.page.tmpl", &models.WebData{MovementSensor1 : broker.WebData.MovementSensor1 , MovementSensor2 : false, LightSensor: broker.WebData.LightSensor, VideosSlice: videos})
	
	} else {
		render.RenderWebPage(w, req, "login.page.tmpl", &models.WebData{})
		fmt.Println("Authentication failure: incorrect username or password")
	}
}


type RegistraionDetails struct {
	Username, Password, PhoneNumber, Email,SlackKey    string 
}


// POST endpoint to process user's registration details 
func (repo *Repository) RegistrationPost(w http.ResponseWriter, req *http.Request) {
 
	Repo.applicationConfiguration.Username = req.Form.Get("fusername")
	Repo.applicationConfiguration.Password = req.Form.Get("fpassword") 
	Repo.applicationConfiguration.Email = req.Form.Get("femail")
	Repo.applicationConfiguration.PhoneNumber = req.Form.Get("fphone")
	Repo.applicationConfiguration.SlackKey = req.Form.Get("fslackKey")

	registraionDetailsData := RegistraionDetails {
		Username: Repo.applicationConfiguration.Username,
		Password: Repo.applicationConfiguration.Password,
		PhoneNumber: Repo.applicationConfiguration.PhoneNumber,
		Email: Repo.applicationConfiguration.Email,
		SlackKey: Repo.applicationConfiguration.SlackKey,
	}

	file, _ := json.MarshalIndent(registraionDetailsData, "", " ")

	error := os.WriteFile(Repo.applicationConfiguration.RegistraionDetailsFile, file, 0644)
	if error != nil {
		log.Fatal("error saving registraionDetails file ", error)
	}

	render.RenderWebPage(w, req, "login.page.tmpl", &models.WebData{})
}