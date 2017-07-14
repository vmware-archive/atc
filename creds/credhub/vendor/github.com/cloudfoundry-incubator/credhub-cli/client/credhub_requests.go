package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"net/url"

	"io"
	"io/ioutil"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func NewSetCertificateRequest(config config.Config, credentialIdentifier string, root string, caName string, cert string, priv string, overwrite bool) *http.Request {
	certificate := models.Certificate{
		Ca:          root,
		Certificate: cert,
		PrivateKey:  priv,
		CaName:      caName,
	}

	return NewSetCredentialRequest(config, "certificate", credentialIdentifier, certificate, overwrite)
}

func NewSetRsaSshRequest(config config.Config, credentialIdentifier, keyType, publicKey, privateKey string, overwrite bool) *http.Request {
	key := models.RsaSsh{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	return NewSetCredentialRequest(config, keyType, credentialIdentifier, key, overwrite)
}

func NewSetUserRequest(config config.Config, credentialIdentifier, username, password string, overwrite bool) *http.Request {
	user := models.User{
		Username: username,
		Password: password,
	}

	return NewSetCredentialRequest(config, "user", credentialIdentifier, user, overwrite)
}

func NewSetRequest(config config.Config, content map[string]interface{}) *http.Request {
	return newCredentialRequest("PUT", config, content)
}

func NewSetCredentialRequest(config config.Config, credentialType string, credentialIdentifier string, content interface{}, overwrite bool) *http.Request {
	credential := models.RequestBody{
		CredentialType: credentialType,
		Name:           credentialIdentifier,
		Value:          content,
		Overwrite:      &overwrite,
	}

	return newCredentialRequest("PUT", config, credential)
}

func NewSetJsonCredentialRequest(config config.Config, credentialType string, credentialIdentifier string, content interface{}, overwrite bool) *http.Request {
	var value interface{}
	valueObject := make(map[string]interface{})
	contentCredential := content.(string)
	err := json.Unmarshal([]byte(contentCredential), &valueObject)

	if err != nil {
		value = content
	} else {
		value = valueObject
	}

	credential := models.RequestBody{
		CredentialType: credentialType,
		Name:           credentialIdentifier,
		Value:          value,
		Overwrite:      &overwrite,
	}

	return newCredentialRequest("PUT", config, credential)
}

func NewGenerateCredentialRequest(config config.Config, identifier string, parameters models.GenerationParameters, value *models.ProvidedValue, credentialType string, overwrite bool) *http.Request {
	generateRequest := models.GenerateRequest{
		Name:           identifier,
		CredentialType: credentialType,
		Overwrite:      &overwrite,
		Parameters:     &parameters,
		Value:          value,
	}

	return newCredentialRequest("POST", config, generateRequest)
}

func NewRegenerateCredentialRequest(config config.Config, identifier string) *http.Request {
	regenerateRequest := models.RegenerateRequest{
		Name:       identifier,
		Regenerate: true,
	}

	return newCredentialRequest("POST", config, regenerateRequest)
}

func NewGetCredentialByNameRequest(config config.Config, name string) *http.Request {
	url := config.ApiURL + "/api/v1/data?name=" + url.QueryEscape(name) + "&current=true"
	return newRequestWithoutBody("GET", config, url)
}

func NewGetCredentialByIdRequest(config config.Config, id string) *http.Request {
	url := config.ApiURL + "/api/v1/data/" + url.QueryEscape(id)

	return newRequestWithoutBody("GET", config, url)
}

func NewDeleteCredentialRequest(config config.Config, identifier string) *http.Request {
	url := config.ApiURL + "/api/v1/data?name=" + url.QueryEscape(identifier)
	return newRequestWithoutBody("DELETE", config, url)
}

func NewInfoRequest(config config.Config) *http.Request {
	url := config.ApiURL + "/info"

	request, _ := http.NewRequest("GET", url, nil)

	return request
}

func NewBodyClone(req *http.Request) io.ReadCloser {
	var result io.ReadCloser = nil
	if req.Body != nil {
		var bodyBytes []byte
		buf := new(bytes.Buffer)
		rc, ok := req.Body.(io.ReadCloser)
		if !ok {
			rc = ioutil.NopCloser(req.Body)
		}
		buf.ReadFrom(rc)
		bodyBytes = buf.Bytes()
		req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		result = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	}
	return result
}

func NewFindAllCredentialPathsRequest(config config.Config) *http.Request {
	url := config.ApiURL + "/api/v1/data?paths=true"
	return newRequestWithoutBody("GET", config, url)
}

func NewFindCredentialsBySubstringRequest(config config.Config, partialIdentifier string) *http.Request {
	urlString := config.ApiURL + "/api/v1/data?name-like=" + url.QueryEscape(partialIdentifier)
	return newRequestWithoutBody("GET", config, urlString)
}

func NewFindCredentialsByPathRequest(config config.Config, path string) *http.Request {
	urlString := config.ApiURL + "/api/v1/data?path=" + url.QueryEscape(path)
	return newRequestWithoutBody("GET", config, urlString)
}

func newCredentialRequest(requestType string, config config.Config, bodyModel interface{}) *http.Request {
	urlString := config.ApiURL + "/api/v1/data"
	return newRequest(requestType, config, urlString, bodyModel)
}

func newRequest(requestType string, config config.Config, url string, bodyModel interface{}) *http.Request {
	var request *http.Request
	body, _ := json.Marshal(bodyModel)
	request, _ = http.NewRequest(requestType, url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.AccessToken)
	return request
}

func newRequestWithoutBody(requestType string, config config.Config, url string) *http.Request {
	var request *http.Request
	request, _ = http.NewRequest(requestType, url, nil)
	request.Header.Set("Authorization", "Bearer "+config.AccessToken)
	return request
}
