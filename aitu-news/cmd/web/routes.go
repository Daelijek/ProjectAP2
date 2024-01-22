package main

import (
	"net/http"
	"sync"
)

type Route func(w http.ResponseWriter, r *http.Request)

var instance *map[string]Route
var once sync.Once

func getRoutes() *map[string]Route {
	once.Do(func() {
		routes := make(map[string]Route)
		routes["/"] = home
		routes["/contacts"] = contacts
		routes["/admin"] = adminPanel
		routes["/article"] = article
		routes["/article/create"] = articleCreate
		routes["/auth/register"] = register
		routes["/auth/login"] = login
		routes["/articles/"] = articles

		instance = &routes
	})
	return instance
}
