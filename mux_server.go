package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"log"
	"fmt"
	"time"
	"encoding/json"
	"os"
)

type Todo struct {
    Name      string `json:"name"`
    Completed bool `json:"completed"`
    Due       time.Time `json:"due"`
}

type Todos []Todo

type RouteHandlerFunc func(w http.ResponseWriter, request *http.Request)


type APIMicroversion struct {
	name string
	route string
	versioned_router *mux.Router
}

func (api_version *APIMicroversion) setup_on_router(router *mux.Router) {
	api_versioned_router := router.PathPrefix(api_version.route).Subrouter()
	api_version.versioned_router = api_versioned_router
}

type RouteHandler struct {
	name string
	route string
	handler_func RouteHandlerFunc
	http_methods []string
}

func (api_version *APIMicroversion) setup_api(handlers []RouteHandler) {
	for _, handler := range handlers {
		api_version.versioned_router.HandleFunc(
			handler.route, handler.handler_func).Methods(handler.http_methods...)
	}
}


func run_http_server(router http.Handler, host string, port int) {
	wrapped_router := handlers.LoggingHandler(os.Stdout, router)
	//http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), wrapped_router)
}


func handle_say_hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the other side of a socket.\n"))
}

/*
Weak point:
 - no clear way to specify acceptable headers regex on router
 */
func simple_server_with_route_handler(
		host string, port int) *mux.Router {

	log.Println(fmt.Sprintf("String web server on %v:%v", host, port))
	router := mux.NewRouter()
	router.Headers("Content-Type: application/(text|json)")
	return router
}

/*
Weak point:
 - no way to use a more than one response handler for a single route
 */
func common_handler(w http.ResponseWriter, request *http.Request) {
	handle_say_hello(w, request)
	route_data_handler(w, request)
}


func route_data_handler(w http.ResponseWriter, request *http.Request) {

	todos := Todos{
        	Todo{Name: "Write presentation"},
        	Todo{Name: "Host meetup"},
    	}
	json.NewEncoder(w).Encode(todos)
}


func main() {
	host, port := "localhost", 8000
	router := simple_server_with_route_handler(host, port)
	versioned_api := APIMicroversion{route: "/v1", name:"v1-api"}
	versioned_api.setup_on_router(router)

	versioned_handlers := []RouteHandler{
		{
			name:"temp",
			route:"/{key}",
			handler_func:common_handler,
			http_methods:[]string{"GET"},
		},
	}

	versioned_api.setup_api(versioned_handlers)
	defer run_http_server(router, host, port)
}
