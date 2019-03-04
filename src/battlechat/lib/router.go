package lib

import (
    "net/http"
    "github.com/gorilla/mux"

    "battlechat/config"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    routes := config.Routes
    for _, route := range routes {
        var handler http.Handler
        handler = route.HandlerFunc
        handler = Logger(handler, route.Name)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)

    }
    return router
}
