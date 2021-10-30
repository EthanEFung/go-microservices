package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)
	l.Println("Welcome!")

	// create handlers
	ho := handlers.NewOK(l)
	hp := handlers.NewProducts(l)

	// create a new server mux and registers the handlers
	sm := http.NewServeMux()
	sm.Handle("/", ho)
	sm.Handle("/products", hp)

	// create a new server
	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
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
