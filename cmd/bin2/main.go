package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	v1 "github.com/axetroy/bin2/v1"
)

func handler(w http.ResponseWriter, r *http.Request) error {
	switch true {
	case strings.HasPrefix(r.URL.Path, "/v1/"):
		return v1.Handle(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
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
