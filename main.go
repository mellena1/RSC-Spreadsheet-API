package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mellena1/RSC-Spreadsheet-API/handler"
)

// RouterCreator a handler that can return a gorilla mux router for path prefixes
type RouterCreator interface {
	Router() *mux.Router
}

func main() {
	router := mux.NewRouter()

	childRouters := getChildRouters()
	for _, c := range childRouters {
		router.PathPrefix(c.PathPrefix).Handler(c.Child.Router())
	}

	http.Handle("/", router)
}

// ChildRouter holds a child handler that can be used for path prefixes
type ChildRouter struct {
	PathPrefix string
	Child      RouterCreator
}

func getChildRouters() []ChildRouter {
	return []ChildRouter{
		{
			PathPrefix: "/standings/",
			Child:      &handler.StandingsHandler{},
		},
	}
}
