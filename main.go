package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest/handlers"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gorHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	// create handlers
	hp := handlers.NewProducts(l)

	// create a new server multiplexer and registers the handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", hp.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", hp.UpdateProducts)
	putRouter.Use(hp.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", hp.AddProduct)
	postRouter.Use(hp.MiddlewareProductValidation)

	// handler for documentation
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir(".")))

	// add cors
	ch := gorHandlers.CORS(gorHandlers.AllowedOrigins([]string{"http://localhost:3000", "https://localhost:3000"}))

	// create a new server
	s := http.Server{
		Addr:         ":8080",
		Handler:      ch(sm),
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		l.Println("Starting server on port :8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received.
	sig := <-c
	l.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
