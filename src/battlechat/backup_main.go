package main

import (
    "log"
    // "net/http"

    "battlechat/lib"
    "battlechat/server"
)

const HTTP_PORT = "8000"
const CHAT_PORT = "4000"

func main() {
    router := lib.NewRouter()
    println(router)

    // log.Fatal(http.ListenAndServe(":" + HTTP_PORT, router))
    go log.Fatal(server.ListenAndServe(CHAT_PORT))
}
