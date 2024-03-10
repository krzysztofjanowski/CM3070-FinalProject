package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)



func TestStarterPageGET(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()
	response, error := webServer.Client().Get(webServer.URL + "/starter")

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestStarterPageGET test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}

func TestLoginPageGET(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	response, error := webServer.Client().Get(webServer.URL + "/login")

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestLoginPageGET test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}

func TestDashboardPageGET(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	response, error := webServer.Client().Get(webServer.URL + "/dashboard")

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestDashboardPageGET test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}

func TestNotificationPageGET(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	response, error := webServer.Client().Get(webServer.URL + "/notifications")

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestNotificationPageGET test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}


func TestPrivacyOnGET(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	response, error := webServer.Client().Get(webServer.URL + "/privacy-on")

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestPrivacyOnGET test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}

func TestPrivacyOFFGET(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	response, error := webServer.Client().Get(webServer.URL + "/privacy-off")

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestPrivacyOFFGET test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}


func TestRecordOnDemandGET(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	response, error := webServer.Client().Get(webServer.URL + "/record-now")

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestRecordOnDemandGET test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}

func TestRegistrationPOST(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	data := url.Values{} 

	// let's add test data 
	data.Add("fusername","szmarcek")
	data.Add("fpassword","123")
	data.Add("femail","ad@no.com")
	data.Add("fphone", "077077077077")
	data.Add("fslackKey", "123459")

	response, error := webServer.Client().PostForm(webServer.URL + "/registration-post", data )

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestRegistrationPOST test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}

func TestLoginPOST(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	data := url.Values{} 

	// let's add test data 
	data.Add("fusername","szmarcek")
	data.Add("fpassword","123")

	response, error := webServer.Client().PostForm(webServer.URL + "/login-post", data )

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestLoginPOST test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}

func TestDashboardPOST(t *testing.T) {
	testHttpHandler := setupHandlers()

	webServer := httptest.NewServer(testHttpHandler )

	defer webServer.Close()

	data := url.Values{} 

	// let's add test data 
	data.Add("ffromDate","2024-01-01 00:00:00")
	data.Add("ftoDate","2024-03-13 00:00:00")
	  
	response, error := webServer.Client().PostForm(webServer.URL + "/dashboard-post", data )

	if error != nil {
		log.Fatal(error)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal("TestDashboardPOST test failed due to incorrect status code it received, should be 200", response.StatusCode)
	}
}