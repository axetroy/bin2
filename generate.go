package bin2

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Script struct {
	ExtName  string
	FileName string
	Content  string
}

func isCurl(useragent string) bool {
	return regexp.MustCompile(`^curl\/`).MatchString(useragent)
}

func isWget(useragent string) bool {
	return regexp.MustCompile(`^Wget\/`).MatchString(useragent)
}

func isPowerShell(useragent string) bool {
	return regexp.MustCompile(`PowerShell\/`).MatchString(useragent)
}

//go:embed install.sh
var shellScript []byte

//go:embed install.ps1
var powerShellScript []byte

func Generate(owner string, repo string, version string, binaryName string, userAgent string) (*Script, error) {
	var (
		err    error
		script []byte
	)

	if version != "" && !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if binaryName == "" {
		binaryName = repo
	}

	if isPowerShell(userAgent) {
		script = powerShellScript
	} else if isCurl(userAgent) || isWget(userAgent) {
		script = shellScript
	} else {
		script = shellScript
	}

	// If no version is specified, the latest version is used
	if version == "" {
		version, err = getLatestRelease(owner, repo)

		if err != nil {
			return nil, err
		}
	}

	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if binaryName == "" {
		binaryName = repo
	}

	t := template.New("install")

	t, err = t.Parse(string(script))

	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer

	err = t.Execute(&tpl, map[string]interface{}{
		"Owner":   owner,
		"Repo":    repo,
		"Version": version,
		"Binary":  binaryName,
	})

	if err != nil {
		return nil, err
	}

	return &Script{
		FileName: "install",
		Content:  tpl.String(),
		ExtName:  path.Ext("install"),
	}, nil
}

func getLatestRelease(owner string, repo string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo))

	if err != nil {
		return "", errors.Wrap(err, "fetch remote version information fail")
	}

	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return "", errors.New(fmt.Sprintf("fetch remote version information and get status code %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", errors.Wrap(err, "read from response body fail")
	}

	type Asset struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}

	type Response struct {
		TagName string  `json:"tag_name"`
		Assets  []Asset `json:"assets"`
	}

	response := Response{}

	if err = json.Unmarshal(body, &response); err != nil {
		return "", errors.Wrap(err, "unmarshal response body fail")
	}

	version := response.TagName

	return version, nil
}
