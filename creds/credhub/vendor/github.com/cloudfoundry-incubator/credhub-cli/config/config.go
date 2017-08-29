package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const AuthClient = "credhub_cli"
const AuthPassword = ""

type Config struct {
	ApiURL             string
	AuthURL            string
	AccessToken        string
	RefreshToken       string
	InsecureSkipVerify bool
	CaCert             []string
}

func ConfigDir() string {
	return path.Join(userHomeDir(), ".credhub")
}

func ConfigPath() string {
	return path.Join(ConfigDir(), "config.json")
}

func ReadConfig() Config {
	c := Config{}

	data, err := ioutil.ReadFile(ConfigPath())
	if err != nil {
		return c
	}

	json.Unmarshal(data, &c)

	return c
}

func WriteConfig(c Config) error {
	err := makeDirectory()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configPath := ConfigPath()
	return ioutil.WriteFile(configPath, data, 0600)
}

func RemoveConfig() error {
	return os.Remove(ConfigPath())
}
