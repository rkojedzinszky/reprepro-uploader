package main

import (
	"context"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/namsral/flag"
	"github.com/rkojedzinszky/reprepro-uploader/pkg/reaper"
	"github.com/rkojedzinszky/reprepro-uploader/pkg/token"
)

var (
	jweSecret     = flag.String("jwe-secret", "", "Base64 encoded JWE token")
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

	secret, err := base64.StdEncoding.DecodeString(*jweSecret)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	reaper := reaper.New()
	wg.Add(1)
	go func() {
		defer wg.Done()

		reaper.Run(ctx)
	}()

	server := &server{
		decoder:      token.MustNewDecoder(secret, token.DecoderWithTime()),
		repreproPath: *repreproPath,
		reaper:       reaper,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		httpserver := &http.Server{}
		mux := http.NewServeMux()
		mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		mux.Handle(*uploadUri, server)
		httpserver.Handler = mux

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
