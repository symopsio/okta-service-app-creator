package main

import (
	"context"
	"fmt"
	"os"

	o "github.com/okta/okta-sdk-golang/okta"
)

func main() {
	// This will load the .okta.yaml file in the current directory
	client, err := o.NewClient(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error creating okta client: %v\n", err)
		return
	}

	config := client.GetConfig()
	fmt.Fprintf(os.Stdout, "Read Client ID: %v\n", config.Okta.Client.ClientId)
}
