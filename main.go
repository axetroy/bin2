package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func isCurl(useragent string) bool {
	return regexp.MustCompile(`^curl\/`).MatchString(useragent)
}

func isWget(useragent string) bool {
	return regexp.MustCompile(`^Wget\/`).MatchString(useragent)
}

func isPowerShell(useragent string) bool {
	return regexp.MustCompile(`PowerShell\/`).MatchString(useragent)
}

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

	filename := "install.sh"

	if isCurl(userAgent) || isWget(userAgent) {
		filename = "install.sh"
	} else if isPowerShell(userAgent) {
		filename = "install.ps1"
	}

	// If no version is specified, the latest version is used
	if version == "" {
		version, err = getLatestRelease(owner, repo)

		if err != nil {
			return err
		}
	}

	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if binaryName == "" {
		binaryName = repo
	}

	t := template.New(filename)

	w.Header().Set("Content-Type", "text/x-shellscript")

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	t, err = t.Parse(string(b))

	if err != nil {
		return err
	}

	err = t.Execute(w, map[string]interface{}{
		"Owner":   owner,
		"Repo":    repo,
		"Version": version,
		"Binary":  binaryName,
	})

	return err
}

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		err := handler(res, req)

		if err != nil {
			panic(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
