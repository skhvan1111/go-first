package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/skhvan1111/go-first/internal/redis"
	"github.com/skhvan1111/go-first/internal/counter"
	"github.com/skhvan1111/go-first/internal/diagnostics"
)

type serverInfo struct {
	port   string
	router http.Handler
	name   string
}

func main() {
	errors := make(chan error, 2)

	redis.Init(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"))

	router := mux.NewRouter()
	router.HandleFunc("/", hello)

	serverConfigurations := []serverInfo{
		{
			port:   getPort("PORT"),
			router: router,
			name:   "application server",
		},
		{
			port:   getPort("DIAG_PORT"),
			router: diagnostics.NewDiagnostics(),
			name:   "diagnostics server",
		},
	}

	servers := make([]*http.Server, 2)

	for i, info := range serverConfigurations {
		go func(info serverInfo, index int) {
			log.Print("The " + info.name + " is preparing to handle connections...")
			servers[index] = &http.Server{
				Addr:    info.port,
				Handler: info.router,
			}
			err := servers[index].ListenAndServe()
			if err != nil {
				errors <- err
			}
		}(info, i)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errors:
		log.Printf("Got an error %v", err)
	case sig := <-interrupt:
		log.Printf("Received the signal %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, s := range servers {
		shutdownError := s.Shutdown(ctx)
		if shutdownError != nil {
			fmt.Println(shutdownError.Error())
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	count, err := counter.Increase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error occured: %v", err)
	} else {
		fmt.Fprintf(w, "Page viewed %v times", count)
	}
}

func getPort(name string) string {
	port := os.Getenv(name)
	if len(port) == 0 {
		log.Fatal(name + " is not set up")
	}
	return ":" + port
}
