package main

import (
  "fmt"
  "time"
  "net/http"
)

type CORSMiddleware struct {
    handler http.Handler
}

func (middleware *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now();

	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://106.14.155.231,http://localhost:5173");

    middleware.handler.ServeHTTP(w, r);

    fmt.Println("%s %s %v", r.Method, r.URL.Path, time.Since(start));
}

func NewCORSMiddleware(handlerToWrap http.Handler) *CORSMiddleware {
    return &CORSMiddleware{handlerToWrap}
}

func main() {
    fileHandler := http.StripPrefix("/fs/", http.FileServer(http.Dir("./static/")));
    http.Handle("/fs/", NewCORSMiddleware(fileHandler))
    err := http.ListenAndServe(":8001", nil)
    fmt.Println(err)
}