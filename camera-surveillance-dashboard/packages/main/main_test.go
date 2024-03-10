package main

import "testing"

func TestInitiateBackend(t *testing.T) {
	error := initiateBackend()
	if error != nil{
		t.Error("Failed to initate the backend.")
	}
}
