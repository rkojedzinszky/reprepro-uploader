package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	httpAuthorization   = "Authorization"
	httpContentEncoding = "Content-Encoding"
)

type server struct {
	token        string
	repreproPath string
}

func extractToken(r *http.Request) string {
	splitted := strings.Split(r.Header.Get(httpAuthorization), " ")

	if len(splitted) != 2 {
		return ""
	}

	if strings.ToLower(splitted[0]) != "bearer" {
		return ""
	}

	return splitted[1]
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.token != "" && extractToken(r) != s.token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd := exec.Command("upload.sh")
	cmd.Dir = s.repreproPath

	if r.Header.Get(httpContentEncoding) == "gzip" {
		gzipreader, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer gzipreader.Close()
		cmd.Stdin = gzipreader
	} else {
		cmd.Stdin = r.Body
	}

	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	dists := r.URL.Query().Get("dists")
	if dists != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("REPREPRO_REPOS=%s", strings.Join(strings.Split(dists, ","), " ")))
	}

	if err := cmd.Run(); err != nil {
		log.Print("E: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out.Bytes())
}
