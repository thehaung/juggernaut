package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thehaung/juggernaut/config"
	"github.com/thehaung/juggernaut/domain"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func NewRoute(router *chi.Mux, conf *config.Config, juggernautUseCase domain.JuggernautUseCase) {
	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	// Logger Middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if conf.HttpServer.ExcludeLogRouterPath[r.URL.Path] {
				// Skip logging for ExcludeLogRouterPath
				next.ServeHTTP(w, r)
				return
			}
			middleware.Logger(next).ServeHTTP(w, r)
		})
	})

	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             4000, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(time.Duration(conf.HttpServer.Timeout) * time.Second))

	// Prometheus' metrics register
	router.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	// Register this router for health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal("Hello World!!!!")
		if err != nil {
			log.Fatal(fmt.Errorf("NewRouter - json.Marshal: %w", err))
		}
		if _, err = w.Write(res); err != nil {
			log.Println(err.Error())
		}
	})

	// Global prefix
	router.Mount("/api", router)
	{
		NewJuggernautHandler(router, juggernautUseCase)
	}
}
