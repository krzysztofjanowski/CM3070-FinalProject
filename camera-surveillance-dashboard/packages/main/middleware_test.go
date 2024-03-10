package main

import (
	"net/http"
	"testing"
)

func TestCrossSiteToken(t *testing.T) {
	var testHandler testingHandler
	
	crossSiteHandler := CrossSiteToken(&testHandler)

	switch checkType := crossSiteHandler.(type) {
		case http.Handler:
			// TestCrossSiteToken correct handler returned 
		default: 
			t.Error("TestCrossSiteToken returned incorrect type, should be an http.Handler but is:", checkType)
	}

}