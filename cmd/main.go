package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/symopsio/okta-service-app-creator/internal"
	"github.com/symopsio/okta-service-app-creator/pkg/config"
)

var appName string
var oktaAPIToken string
var oktaOrgName string
var outputDirectory string
var outputFile string

func init() {
	flag.StringVar(&appName, "name", "", "Name of the Okta App to create")
	flag.StringVar(&oktaOrgName, "org", "", "Name of the Okta Org (like dev-12345678)")
	flag.StringVar(&oktaAPIToken, "token", "", "Okta API token (can also be set using OKTA_API_TOKEN env var)")
	flag.StringVar(&outputDirectory, "output", "", "Where to write the Okta App private key (default to current directory)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `%s:

Creates an Okta OAuth2 service app using a JSON Web Key (JWK) for client credentials. Writes the configuration to a yaml file in the given output directory, using the YAML format supported by the Okta SDK for golang.

More info: https://github.com/symopsio/okta-service-app-creator

`, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	parseFlags()

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	clientID, err := internal.CreateOktaApp(appName, privateKey,
		oktaOrgName, oktaAPIToken)
	exitOnError(err)

	orgURL := fmt.Sprintf("https://%s.okta.com", oktaOrgName)
	oauthConfig := config.NewConfig(privateKey, clientID, orgURL)
	err = config.WriteConfigFile(oauthConfig, outputFile)
	exitOnError(err)

	fmt.Fprintf(os.Stdout, "Created Okta App with Client ID: %v\n", clientID)
}

func parseFlags() {
	flag.Parse()
	if oktaAPIToken == "" {
		oktaAPIToken = os.Getenv("OKTA_API_TOKEN")
		if oktaAPIToken == "" {
			exitWithUsage("OKTA_API_TOKEN env var is required")
		}
	}
	if appName == "" {
		exitWithUsage("-name is required")
	}
	if oktaOrgName == "" {
		exitWithUsage("-org is required")
	}
	if outputDirectory == "" {
		outputFile = fmt.Sprintf("%s-%s.yaml", oktaOrgName, appName)
	} else {
		outputFile = fmt.Sprintf("%s/%s-%s.yaml", oktaOrgName, outputDirectory, appName)
	}
	if !isValidFile(outputFile) {
		exitOnError(fmt.Errorf("Invalid file path: %s", outputFile))
	}
}

func isValidFile(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}

func exitWithUsage(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	flag.Usage()
	os.Exit(1)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
