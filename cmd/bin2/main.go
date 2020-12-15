package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/axetroy/bin2"
)

func getLatestRelease(owner string, repo string) (string, error) {
	var (
		err error
	)

	url := fmt.Sprintf("https://github.com/%s/%s/releases", owner, repo)

	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(fmt.Sprintf(`/%s/%s/releases/download/(.*)/.*\.tar.gz`, owner, repo))
	// reg := regexp.MustCompile(fmt.Sprintf(`/%s/%s/releases/tag/v([a-z\d\.]+)`, owner, repo))

	matchers := reg.FindStringSubmatch(string(body))

	if len(matchers) == 0 {
		return "", errors.New("not found latest version")
	}

	version := matchers[1]

	return version, nil
}

func handler(w http.ResponseWriter, r *http.Request) error {
	var (
		err        error
		version    = r.URL.Query().Get("v")
		binaryName = r.URL.Query().Get("bin")
	)

	arr := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(arr) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	owner := arr[0]
	repo := strings.Join(arr[1:], "/")

	userAgent := r.Header.Get("user-agent")

	// If no version is specified, the latest version is used
	if version == "" {
		version, err = getLatestRelease(owner, repo)

		if err != nil {
			return err
		}
	}

	script, err := bin2.Generate(owner, repo, version, binaryName, userAgent)

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
