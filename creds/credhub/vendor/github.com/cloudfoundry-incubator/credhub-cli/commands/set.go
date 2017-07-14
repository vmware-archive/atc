package commands

import (
	"fmt"

	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"

	"bufio"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type SetCommand struct {
	CredentialIdentifier string `short:"n" required:"yes" long:"name" description:"Name of the credential to set"`
	Type                 string `short:"t" long:"type" description:"Sets the credential type. Valid types include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'. Type-specific values are set with the following flags (supported types prefixed)."`
	NoOverwrite          bool   `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
	Value                string `short:"v" long:"value" description:"[Value, JSON] Sets the value for the credential"`
	CaName               string `short:"m" long:"ca-name" description:"[Certificate] Sets the root CA to a stored CA credential"`
	Root                 string `short:"r" long:"root" description:"[Certificate] Sets the root CA from file"`
	Certificate          string `short:"c" long:"certificate" description:"[Certificate] Sets the certificate from file"`
	Private              string `short:"p" long:"private" description:"[Certificate, SSH, RSA] Sets the private key from file"`
	Public               string `short:"u" long:"public" description:"[SSH, RSA] Sets the public key from file"`
	RootString           string `short:"R" long:"root-string" description:"[Certificate] Sets the root CA from string input"`
	CertificateString    string `short:"C" long:"certificate-string" description:"[Certificate] Sets the certificate from string input"`
	PrivateString        string `short:"P" long:"private-string" description:"[Certificate, SSH, RSA] Sets the private key from string input"`
	PublicString         string `short:"U" long:"public-string" description:"[SSH, RSA] Sets the public key from  string input"`
	Username             string `short:"z" long:"username" description:"[User] Sets the username value of the credential"`
	Password             string `short:"w" long:"password" description:"[Password, User] Sets the password value of the credential"`
	OutputJson           bool   `          long:"output-json" description:"Return response in JSON format"`
}

func (cmd SetCommand) Execute([]string) error {

	if cmd.Type == "" {
		return errors.NewSetEmptyTypeError()
	}

	if cmd.Value == "" && (cmd.Type == "value" || cmd.Type == "json") {
		promptForInput("value: ", &cmd.Value)
	}

	if cmd.Password == "" && (cmd.Type == "password" || cmd.Type == "user") {
		promptForInput("password: ", &cmd.Password)
	}

	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))

	action := actions.NewAction(repository, &cfg)
	request, err := MakeRequest(cmd, cfg)
	if err != nil {
		return err
	}

	credential, err := action.DoAction(request, cmd.CredentialIdentifier)
	if err != nil {
		return err
	}
	models.Println(credential, cmd.OutputJson)

	return nil
}

func MakeRequest(cmd SetCommand, config config.Config) (*http.Request, error) {
	var request *http.Request
	if cmd.Type == "ssh" || cmd.Type == "rsa" {
		var err error

		err = setStringFieldFromFile(&cmd.Public, &cmd.PublicString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&cmd.Private, &cmd.PrivateString)
		if err != nil {
			return nil, err
		}

		request = client.NewSetRsaSshRequest(config, cmd.CredentialIdentifier, cmd.Type, cmd.PublicString, cmd.PrivateString, !cmd.NoOverwrite)
	} else if cmd.Type == "certificate" {
		var err error

		err = setStringFieldFromFile(&cmd.Root, &cmd.RootString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&cmd.Certificate, &cmd.CertificateString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&cmd.Private, &cmd.PrivateString)
		if err != nil {
			return nil, err
		}

		request = client.NewSetCertificateRequest(config, cmd.CredentialIdentifier, cmd.RootString, cmd.CaName, cmd.CertificateString, cmd.PrivateString, !cmd.NoOverwrite)
	} else if cmd.Type == "user" {
		request = client.NewSetUserRequest(config, cmd.CredentialIdentifier, cmd.Username, cmd.Password, !cmd.NoOverwrite)
	} else if cmd.Type == "password" {
		request = client.NewSetCredentialRequest(config, cmd.Type, cmd.CredentialIdentifier, cmd.Password, !cmd.NoOverwrite)
	} else if cmd.Type == "json" {
		request = client.NewSetJsonCredentialRequest(config, cmd.Type, cmd.CredentialIdentifier, cmd.Value, !cmd.NoOverwrite)
	} else {
		request = client.NewSetCredentialRequest(config, cmd.Type, cmd.CredentialIdentifier, cmd.Value, !cmd.NoOverwrite)
	}

	return request, nil
}

func promptForInput(prompt string, value *string) {
	fmt.Printf(prompt)
	reader := bufio.NewReader(os.Stdin)
	val, _ := reader.ReadString('\n')
	*value = string(strings.TrimSpace(val))
}

func setStringFieldFromFile(fileField, stringField *string) error {
	var err error
	if *fileField != "" {
		if *stringField != "" {
			return errors.NewCombinationOfParametersError()
		}
		*stringField, err = ReadFile(*fileField)
		if err != nil {
			return err
		}
	}
	return nil
}
