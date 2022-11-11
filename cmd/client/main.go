package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/namsral/flag"
)

var (
	outputPath     = flag.String("reprepro-upload-path", "/output", "Path for .deb files")
	repreproToken  = flag.String("reprepro-token", "", "Token")
	repreproServer = flag.String("reprepro-server", "", "Server address")
	distributions  = flag.String("distribution", "", "Distributions, separated by comma")
)

func main() {
	flag.Parse()

	if err := os.Chdir(*outputPath); err != nil {
		log.Fatal(err)
	}

	r, w := io.Pipe()

	go createTar(w)

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s?dists=%s", *repreproServer, *distributions), r)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %s", *repreproToken))

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		log.Fatal("Expected http accepted, got: ", resp.StatusCode)
	}

	io.Copy(os.Stdout, resp.Body)
}

func createTar(w io.WriteCloser) {
	defer w.Close()
	gz := gzip.NewWriter(w)
	defer gz.Close()
	tarf := tar.NewWriter(gz)
	defer tarf.Close()

	files, _ := filepath.Glob("*")
	for _, file := range files {
		fstat, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}

		if fstat.Mode()&os.ModeType != 0 {
			continue
		}

		if err := tarf.WriteHeader(&tar.Header{
			Name: fstat.Name(),
			Size: fstat.Size(),
			Mode: 0o644,
		}); err != nil {
			log.Fatal(err)
		}

		fh, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.Copy(tarf, fh); err != nil {
			log.Fatal(err)
		}
		fh.Close()
	}
}
