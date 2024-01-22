package main

import (
	"net/http"
	"sync"
)

type Route func(w http.ResponseWriter, r *http.Request)

var (
	instance *map[string]Route
	once     sync.Once
)

func getRoutes() *map[string]Route {
	once.Do(func() {
		routes := map[string]Route{
			"/":               home,
			"/contacts":       contacts,
			"/admin":          adminPanel,
			"/article":        article,
			"/article/create": articleCreate,
			"/auth/register":  register,
			"/auth/login":     login,
			"/articles/":      articles,
		}

		instance = &routes
	})
	return instance
}

// Define your route handle
