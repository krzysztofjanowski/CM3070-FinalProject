package main

import (
	"testing"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/config"
	"net/http"
)

func TestRoutes(t *testing.T){

	var appConfig config.AppConfig;

	router := routes(&appConfig) 

	switch checkType := router.(type) {
		case http.Handler:
			// TestCrossSiteToken correct handler returned 
		default: 
			t.Error("TestRoutes returned incorrect type, should be an http.Handler but is:", checkType)
	} 

}