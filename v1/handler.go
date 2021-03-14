package v1

import (
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) error {
	var (
		err        error
		version    = r.URL.Query().Get("v")
		binaryName = r.URL.Query().Get("bin")
		binDir     = r.URL.Query().Get("dir")
	)

	urlPath := strings.Trim(r.URL.Path, "/v1/")

	arr := strings.Split(urlPath, "/")

	if len(arr) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	owner := arr[0]
	repo := strings.Join(arr[1:], "/")
	userAgent := r.Header.Get("user-agent")

	script, err := generate(owner, repo, version, binaryName, binDir, userAgent)

	if err != nil {
		return err
	}

	_, err = w.Write([]byte(script.Content))

	return err
}
