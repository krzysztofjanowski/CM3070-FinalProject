package main

import (
	"net/http"
	"github.com/justinas/nosurf"
)

//  csrf protection middleware necessary for forms 
func CrossSiteToken(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	return csrfHandler
}