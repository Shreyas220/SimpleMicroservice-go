package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Shreyas220/SimpleMicroservice-go/handlers"
)

func main() {

	//convience method (default server munx)
	l := log.New(os.Stdout, "SimpleMicroservice-go", log.LstdFlags)
	hh := handlers.NewHello(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	//converting a function into a handler
	//http.HandleFunc()

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
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
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
