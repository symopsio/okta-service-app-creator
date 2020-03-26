package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Okta Okta `yaml:"okta"`
}

type Okta struct {
	Client Client `yaml:"client"`
}

type Client struct {
	OrgUrl            string `yaml:"orgUrl"`
	AuthorizationMode string `yaml:"authorizationMode"`
	ClientId          string `yaml:"clientId"`
	PrivateKey        string `yaml:"privateKey"`
}

// NewConfig creates a config by encoding the rsa private key as a pem key
func NewConfig(privateKey *rsa.PrivateKey, clientID, orgURL string) *Config {
	pemKey := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	return &Config{
		Okta{
			Client{
				AuthorizationMode: "PrivateKey",
				OrgUrl:            orgURL,
				ClientId:          clientID,
				PrivateKey:        string(pemKey),
			},
		},
	}
}

// WriteConfigFile writes a YAML version of the config to the given path
func WriteConfigFile(config *Config, filePath string) error {
	dat, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filePath, dat, 0644); err != nil {
		return fmt.Errorf("Unable to write to file %q: %w", filePath, err)
	}
	fmt.Fprintf(os.Stdout, "Wrote config to file: %s\n", filePath)
	return nil
}

// ReadConfigFile reads a YAML version of the config from the given path
func ReadConfigFile(filePath string) (*Config, error) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read from file %q: %w", filePath, err)
	}
	config := &Config{}
	if err = yaml.Unmarshal(dat, &config); err != nil {
		return nil, fmt.Errorf("Unable to read from file %q: %w", filePath, err)
	}
	return config, nil
}
