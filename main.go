package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"

	"yss-go-official/logger"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(logger.LoggerMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		logger.GetLogEntry(r).Info("hello")
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error at listen", err)
	}
}
