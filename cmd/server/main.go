package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/namsral/flag"
)

var (
	token         = flag.String("token", "", "Bearer token used for authenticating clients")
	uploadUri     = flag.String("upload-uri", "/upload", "Upload uri")
	repreproPath  = flag.String("reprepro-path", "/home/reprepro", "Path to reprepro home")
	listenAddress = flag.String("listen-address", ":8080", "Listen address")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &server{
		token:        *token,
		repreproPath: *repreproPath,
	}

	httpserver := &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle(*uploadUri, server)
	httpserver.Handler = mux

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		go func() {
			<-ctx.Done()

			httpserver.Close()
		}()

		httpserver.Serve(lis)
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)
	<-sigchan

	log.Print("Exiting")

	cancel()

	wg.Wait()
}
