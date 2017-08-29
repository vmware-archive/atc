package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type GenerateCommand struct {
	CredentialIdentifier string   `short:"n" required:"yes" long:"name" description:"Name of the credential to generate"`
	CredentialType       string   `short:"t" long:"type" description:"Sets the credential type to generate. Valid types include 'password', 'user', 'certificate', 'ssh' and 'rsa'. Type-specific values are set with the following flags (supported types prefixed)."`
	NoOverwrite          bool     `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
	OutputJson           bool     `long:"output-json" description:"Return response in JSON format"`
	Username             string   `short:"z" long:"username" description:"Sets the username value of the credential"`
	Length               int      `short:"l" long:"length" description:"[Password, User] Length of the generated value (Default: 30)"`
	IncludeSpecial       bool     `short:"S" long:"include-special" description:"[Password, User] Include special characters in the generated value"`
	ExcludeNumber        bool     `short:"N" long:"exclude-number" description:"[Password, User] Exclude number characters from the generated value"`
	ExcludeUpper         bool     `short:"U" long:"exclude-upper" description:"[Password, User] Exclude upper alpha characters from the generated value"`
	ExcludeLower         bool     `short:"L" long:"exclude-lower" description:"[Password, User] Exclude lower alpha characters from the generated value"`
	SshComment           string   `short:"m" long:"ssh-comment" description:"[SSH] Comment appended to public key to help identify in environment"`
	KeyLength            int      `short:"k" long:"key-length" description:"[Certificate, SSH, RSA] Bit length of the generated key (Default: 2048)"`
	Duration             int      `short:"d" long:"duration" description:"[Certificate] Valid duration (in days) of the generated certificate (Default: 365)"`
	CommonName           string   `short:"c" long:"common-name" description:"[Certificate] Common name of the generated certificate"`
	Organization         string   `short:"o" long:"organization" description:"[Certificate] Organization of the generated certificate"`
	OrganizationUnit     string   `short:"u" long:"organization-unit" description:"[Certificate] Organization unit of the generated certificate"`
	Locality             string   `short:"i" long:"locality" description:"[Certificate] Locality/city of the generated certificate"`
	State                string   `short:"s" long:"state" description:"[Certificate] State/province of the generated certificate"`
	Country              string   `short:"y" long:"country" description:"[Certificate] Country of the generated certificate"`
	AlternativeName      []string `short:"a" long:"alternative-name" description:"[Certificate] A subject alternative name of the generated certificate (may be specified multiple times)"`
	KeyUsage             []string `short:"g" long:"key-usage" description:"[Certificate] Key Usage extensions for the generated certificate (may be specified multiple times)"`
	ExtendedKeyUsage     []string `short:"e" long:"ext-key-usage" description:"[Certificate] Extended Key Usage extensions for the generated certificate (may be specified multiple times)"`
	Ca                   string   `long:"ca" description:"[Certificate] Name of CA used to sign the generated certificate"`
	IsCA                 bool     `long:"is-ca" description:"[Certificate] The generated certificate is a certificate authority"`
	SelfSign             bool     `long:"self-sign" description:"[Certificate] The generated certificate will be self-signed"`
}

func (cmd GenerateCommand) Execute([]string) error {
	if cmd.CredentialType == "" {
		return errors.NewGenerateEmptyTypeError()
	}

	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))

	parameters := models.GenerationParameters{
		IncludeSpecial:   cmd.IncludeSpecial,
		ExcludeNumber:    cmd.ExcludeNumber,
		ExcludeUpper:     cmd.ExcludeUpper,
		ExcludeLower:     cmd.ExcludeLower,
		Length:           cmd.Length,
		CommonName:       cmd.CommonName,
		Organization:     cmd.Organization,
		OrganizationUnit: cmd.OrganizationUnit,
		Locality:         cmd.Locality,
		State:            cmd.State,
		Country:          cmd.Country,
		AlternativeName:  cmd.AlternativeName,
		ExtendedKeyUsage: cmd.ExtendedKeyUsage,
		KeyUsage:         cmd.KeyUsage,
		KeyLength:        cmd.KeyLength,
		Duration:         cmd.Duration,
		Ca:               cmd.Ca,
		SelfSign:         cmd.SelfSign,
		IsCA:             cmd.IsCA,
		SshComment:       cmd.SshComment,
	}

	var value *models.ProvidedValue
	if len(cmd.Username) > 0 {
		value = &models.ProvidedValue{
			Username: cmd.Username,
		}
	}

	action := actions.NewAction(repository, &cfg)
	request := client.NewGenerateCredentialRequest(cfg, cmd.CredentialIdentifier, parameters, value, cmd.CredentialType, !cmd.NoOverwrite)
	credential, err := action.DoAction(request, cmd.CredentialIdentifier)

	if err != nil {
		return err
	}

	models.Println(credential, cmd.OutputJson)

	return nil
}
