package bin2

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
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

func Generate(owner string, repo string, version string, binaryName string, userAgent string) (*Script, error) {
	var (
		err      error
		filename = "install.sh"
	)

	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if binaryName == "" {
		binaryName = repo
	}

	if isCurl(userAgent) || isWget(userAgent) {
		filename = "install.sh"
	} else if isPowerShell(userAgent) {
		filename = "install.ps1"
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

	t := template.New(filename)

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	t, err = t.Parse(string(b))

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
		FileName: filename,
		Content:  tpl.String(),
		ExtName:  path.Ext(filename),
	}, nil
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
