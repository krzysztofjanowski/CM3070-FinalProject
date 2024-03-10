package config

import (
	"html/template"
	"github.com/alexedwards/scs/v2"
)

// main app configurtion holders 
type AppConfig struct {
	TemplateRootDirectory string 
	TemplateCache map[string]*template.Template
	NotificationsLog string 
	VideoDir 	string 
	Username    string 
	Password    string 
	PhoneNumber string 
	Email		string
	SlackKey    string 
	Session  	*scs.SessionManager
	RegistraionDetailsFile string 
	UserAlreadyRegistered bool 
}