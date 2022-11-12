package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/namsral/flag"

	"github.com/rkojedzinszky/reprepro-uploader/pkg/claims"
	"github.com/rkojedzinszky/reprepro-uploader/pkg/token"
)

var (
	outputPath     = flag.String("reprepro-upload-path", "/output", "Path for .deb files")
	jweSecret      = flag.String("jwe-secret", "", "Base64 encoded JWE token")
	repreproServer = flag.String("reprepro-server", "", "Server address")
	distributions  = flag.String("distribution", "", "Distributions, separated by comma")
)

func main() {
	flag.Parse()

	if err := os.Chdir(*outputPath); err != nil {
		log.Fatal(err)
	}

	token := genToken()

	r, w := io.Pipe()

	go createTar(w)

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s?dists=%s", *repreproServer, *distributions), r)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		fmt.Println("Unexpected response: ", resp.StatusCode)
	} else {
		fmt.Println("Success")
	}

	fmt.Print("\nResponse:\n\n")

	io.Copy(os.Stdout, resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		os.Exit(1)
	}
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

func genToken() string {
	secret, err := base64.StdEncoding.DecodeString(*jweSecret)
	if err != nil {
		log.Fatal(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	encoder, err := token.NewEncoder(secret, token.EncoderWithAge(5*time.Second), token.EncoderWithIssuer(hostname))
	if err != nil {
		log.Fatal(err)
	}

	claim := &claims.Claims{
		Distributions: strings.Split(*distributions, ","),
	}

	encoded, err := encoder.Encode(claim)
	if err != nil {
		log.Fatal(err)
	}

	return encoded
}
