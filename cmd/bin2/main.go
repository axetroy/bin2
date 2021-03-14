package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/axetroy/bin2"
)

func handler(w http.ResponseWriter, r *http.Request) error {
	var (
		err        error
		version    = r.URL.Query().Get("v")
		binaryName = r.URL.Query().Get("bin")
		binDir     = r.URL.Query().Get("dir")
	)

	arr := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(arr) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	owner := arr[0]
	repo := strings.Join(arr[1:], "/")
	userAgent := r.Header.Get("user-agent")

	script, err := bin2.Generate(owner, repo, version, binaryName, binDir, userAgent)

	if err != nil {
		return err
	}

	_, err = w.Write([]byte(script.Content))

	return err
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		err := handler(res, req)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, _ = res.Write([]byte(err.Error()))
		}
	})

	fmt.Printf("Listen on port %s\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
