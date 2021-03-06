package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Shreyas220/SimpleMicroservice-go/handlers"
	"github.com/gorilla/mux"
)

func main() {

	//convience method (default server munx)
	l := log.New(os.Stdout, "SimpleMicroservice-go", log.LstdFlags)
	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	//-> middleware -> Update
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProducts)
	postRouter.Use(ph.MiddlewareValidateProduct)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	//contructs an http server
	//http.ListenAndServe(":9090", sm)

	go func() {
		l.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Println("REcieved terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}

//when a request comes there is default handler --> http ServeMux it determines which handler registered against to call  and pass through response writer and request
//Time Out - resources are Finite....multiple of block connections --> basic denial of service attack
//graceful ShutDown (want to shut down Server for any reason ) s.Shutdown --> will wait for all request to complete and then shutdown(no new request will be accepted)
