package config

import (
    "net/http"
    "github.com/gorilla/mux"

    "battlechat/app/controllers"
)

type Route struct {
    Name string
    Method string
    Pattern string
    HandlerFunc http.HandlerFunc
}

type routesArr []Route

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range Routes {
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(route.HandlerFunc)
    }
    return router
}

var Routes = routesArr {
    Route { "Index", "GET", "/", controllers.Index, },
    Route { "TodoIndex", "GET", "/todos", controllers.TodoIndex,},
    Route { "TodoShow", "GET", "/todos/{todoId}", controllers.TodoShow, },
}
