package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/bounce", bounceHandler)
	r.HandleFunc("/ui", uiHandler)

	srv := &http.Server{
		Addr:    ":8001",
		Handler: r,
	}

	log.Fatal(srv.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	_, _ = w.Write([]byte(`<a href="/bounce">Bounce!</a>`))
}

func bounceHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue := fmt.Sprintf("cookie-time-%d", time.Now().UTC().UnixMilli())
	w.Header().Add("Set-Cookie", fmt.Sprintf("dummy-cookie=%s; SameSite=Strict; Max-Age=120", cookieValue))
	w.Header().Set("Location", "/ui")

	w.WriteHeader(http.StatusTemporaryRedirect)
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Header.Values("cookie")
	_, _ = w.Write([]byte("Found these cookies:\n"))
	for _, cookie := range cookies {
		_, _ = w.Write([]byte("   " + cookie + "\n"))
	}
}
