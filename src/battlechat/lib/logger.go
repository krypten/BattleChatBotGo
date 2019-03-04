package lib

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        inner.ServeHTTP(w, r)

        log.Printf(
            "%s\t%s\t%s\t%s",
            r.Method,
            r.RequestURI,
            name,
            time.Since(start),
        )
    })
}

func LoggerE(err error, format string, args ...interface{}) {
    LoggerF(format, args);
    fmt.Println("Error : ", err.Error())
}

func LoggerF(format string, args ...interface{}) {
    log.Printf(format, args);
}
