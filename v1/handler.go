package v1

import (
	"net/http"
	"strings"

	"github.com/axetroy/bin2"
)

func Handle(w http.ResponseWriter, r *http.Request) error {
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
