package render

import (
	"fmt"
	"html/template"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/config"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/models"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

var applicationConfiguration *config.AppConfig

// obtain the application configuration data which contain templates and other  configuration pieces
func PassConfig(a *config.AppConfig) {
	applicationConfiguration = a
}

 
// allows adding some default data before passing it to the actual template for rending
func SetDefaultData(WebData *models.WebData, req *http.Request) *models.WebData {
	WebData.CrossSiteToken = nosurf.Token(req)
	return WebData
}

 
func RenderWebPage(w http.ResponseWriter, req *http.Request, templateName string, WebData *models.WebData) {

	// new templates or layouts need to be added here 
	parsedTemplate, err := template.ParseFiles(fmt.Sprintf("%s/web/templates/", applicationConfiguration.TemplateRootDirectory)+templateName,
											   fmt.Sprintf("%s/web/templates/base.layout.tmpl", applicationConfiguration.TemplateRootDirectory),
											   fmt.Sprintf("%s/web/templates/navigation.layout.tmpl", applicationConfiguration.TemplateRootDirectory))
	if err != nil {
		log.Fatal("error obtaing template: ", err)
	}

	WebData = SetDefaultData(WebData, req)

	// pass the data to the template and render it 
	err = parsedTemplate.Execute(w, WebData)
	if err != nil {
		fmt.Println("error executing tempalte", err)
		return
	}
}

