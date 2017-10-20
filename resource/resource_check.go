package resource

import (
	"github.com/concourse/atc"
	"github.com/tedsuo/ifrit"
)

type checkRequest struct {
	Source  atc.Source  `json:"source"`
	Space   string      `json:"space"`
	Version atc.Version `json:"version"`
}

func (resource *resource) Check(source atc.Source, space string, fromVersion atc.Version) ([]atc.Version, error) {
	var versions []atc.Version

	checking := ifrit.Invoke(resource.runScript(
		"/opt/resource/check",
		nil,
		checkRequest{source, space, fromVersion},
		&versions,
		nil,
		false,
	))

	err := <-checking.Wait()
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (resource *resource) CheckSpaces(source atc.Source) ([]string, error) {
	var spaces []string

	checking := ifrit.Invoke(resource.runScript(
		"/opt/resource/spaces",
		nil,
		checkRequest{Source: source},
		&spaces,
		nil,
		false,
	))

	err := <-checking.Wait()
	if err != nil {
		return nil, err
	}

	return spaces, nil
}
