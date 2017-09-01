package credhub

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// Request sends an authenticated request to the CredHub server.
//
// The pathStr should include the full path (eg. /api/v1/data).
// The request body should be marshallable to JSON, but can be left nil for GET requests.
//
// Request() is used by other CredHub client methods to send authenticated requests to the CredHub server.
//
// Use Request() directly to send authenticated requests to the CredHub server.
// For unauthenticated requests (eg. /health), use Config.Client() instead.
func (ch *CredHub) Request(method string, pathStr string, query url.Values, body interface{}) (*http.Response, error) {
	return ch.request(ch.Auth, method, pathStr, query, body)
}

type requester interface {
	Do(req *http.Request) (*http.Response, error)
}

func (ch *CredHub) request(client requester, method string, pathStr string, query url.Values, body interface{}) (*http.Response, error) {
	u := *ch.baseURL // clone
	u.Path = pathStr
	u.RawQuery = query.Encode()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return resp, err
	}

	if err := ch.checkForServerError(resp); err != nil {
		return nil, err
	}

	return resp, err
}

func (ch *CredHub) checkForServerError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)

		respErr := &Error{}

		if err := dec.Decode(respErr); err != nil {
			return err
		}

		return respErr
	}

	return nil
}
