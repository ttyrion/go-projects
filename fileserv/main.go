package main

import (
  "fmt"
  "net/http"
)

type CORSMiddleware struct {
    handler http.Handler;
    allowedOrigins []string;
}

func NewCORSMiddleware(handlerToWrap http.Handler, origins []string) *CORSMiddleware {
    return &CORSMiddleware{handlerToWrap, origins};
}

func (middleware *CORSMiddleware) IsAllowedOriginsSet() bool {
    return len(middleware.allowedOrigins) > 0;
}

func (middleware *CORSMiddleware) isAllowedOriginsMatched(origin string) bool {
    if middleware.IsAllowedOriginsSet() {
        for _, o := range middleware.allowedOrigins {
            if o == origin {
                return true;
            }
        }
    }

    return false;
}

func (middleware *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS
    if len(r.Header["Origin"]) > 0 {
        requestOrigin := r.Header["Origin"][0];
        // origin matched
        if middleware.isAllowedOriginsMatched(requestOrigin) {
	        w.Header().Set("Access-Control-Allow-Origin", requestOrigin);
        }
    }

    middleware.handler.ServeHTTP(w, r);
}


func main() {
    fileHandler := http.StripPrefix("/fs/", http.FileServer(http.Dir("./static/")));

    allowedOrigins := []string{
        "http://www.huicentury.tech",
        "http://106.14.155.231",
        "http://localhost:5173",
        "http://192.168.1.9:3000",
    }
    http.Handle("/fs/", NewCORSMiddleware(fileHandler, allowedOrigins))

    err := http.ListenAndServe(":8001", nil)
    fmt.Println(err)
}