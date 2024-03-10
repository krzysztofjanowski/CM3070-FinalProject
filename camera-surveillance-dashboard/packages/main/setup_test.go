package main

import (
	"net/http"
	"os"
	"testing"
)

// we need an http handler needed for some testing such as middleware 
type testingHandler struct {

}

// we need this to make testHandler satisfy http.Handler type requiremetns 
func (th *testingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
}

// initates objects needed for running tests
func TestMain(m *testing.M){

	// everything that needs to be set up before running tests comes here


	// run the actual tests 
	os.Exit(m.Run())

}