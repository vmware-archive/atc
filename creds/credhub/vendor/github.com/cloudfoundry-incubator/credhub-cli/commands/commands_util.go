package commands

import (
	"io/ioutil"

	"strings"

	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

func ReadFile(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", credhub_errors.NewFileLoadError()
	}
	return string(dat), nil
}

func AddDefaultSchemeIfNecessary(serverUrl string) string {
	if strings.Contains(serverUrl, "://") {
		return serverUrl
	} else {
		return "https://" + serverUrl
	}
}
