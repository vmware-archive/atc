package errors

import (
	"errors"
	"fmt"
)

func NewNetworkError(e error) error {
	return errors.New(fmt.Sprintf("Error connecting to the targeted API: %#v. Please validate your target and retry your request.", e.Error()))
}

func NewResponseError() error {
	return errors.New("An error occurred when processing the response. Please validate your input and retry your request.")
}

func NewCatchAllError() error {
	return errors.New("The targeted API was unable to perform the request. Please validate and retry your request.")
}

func NewInvalidTargetError() error {
	return errors.New("The targeted API does not appear to be valid. Please validate the API address and retry your request.")
}

func NewNoTargetUrlError() error {
	return errors.New("An API target is not set. Please target the location of your server with `credhub api --server api.example.com` to continue.")
}

func NewRevokedTokenError() error {
	return errors.New("You are not currently authenticated. Please log in to continue.")
}

func NewFileLoadError() error {
	return errors.New("A referenced file could not be opened. Please validate the provided filenames and permissions, then retry your request.")
}

func NewMissingGetParametersError() error {
	return errors.New("A name or ID must be provided. Please update and retry your request.")
}

func NewCombinationOfParametersError() error {
	return errors.New("The combination of parameters in the request is not allowed. Please validate your input and retry your request.")
}

func NewAuthorizationError() error {
	return errors.New("The provided username and password combination are incorrect. Please validate your input and retry your request.")
}

func NewMixedAuthorizationParametersError() error {
	return errors.New("Client and password credentials may not be combined. Please update and retry your request with a single login method.")
}

func NewPasswordAuthorizationParametersError() error {
	return errors.New("The combination of parameters in the request is not allowed. Please validate your input and retry your request.")
}

func NewClientAuthorizationParametersError() error {
	return errors.New("Both client name and client secret must be provided to authenticate. Please update and retry your request.")
}

func NewRefreshError() error {
	return errors.New("You are not currently authenticated. Please log in to continue.")
}

func NewForbiddenError() error {
	return errors.New("You are not authorized to perform this action. You must log in with an elevated user or contact your administrator to continue.")
}

func NewNoMatchingCredentialsFoundError() error {
	return errors.New("No credentials exist which match the provided parameters.")
}

func NewAccessTokenExpiredError() error {
	return errors.New("JWT access token expired")
}

func NewSetEmptyTypeError() error {
	return errors.New("A type must be specified when setting a credential. Valid types include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'.")
}

func NewGenerateEmptyTypeError() error {
	return errors.New("A type must be specified when generating a credential. Valid types include 'password', 'user', 'certificate', 'ssh' and 'rsa'.")
}

func NewNoApiUrlSetError() error {
	return errors.New("An API target is not set. Please target the location of your server with `credhub api --server api.example.com` to continue.")
}

func NewInvalidImportYamlError() error {
	return errors.New("The referenced file does not contain valid yaml structure. Please update and retry your request.")
}
