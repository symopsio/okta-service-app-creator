package internal

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
)

// CreateOktaApp creates an oauth2 service-to-service app.
// Based on
// https://developer.okta.com/docs/reference/api/apps/#add-oauth-2-0-client-application
//
// The okta golang sdk structs do not support all the properties we need,
// so we're doing our own HTTP processing instead.
func CreateOktaApp(appName string, privateKey *rsa.PrivateKey, orgName, token string) (string, error) {
	requestURL := fmt.Sprintf("https://%s.okta.com/api/v1/apps", orgName)

	payload, err := createOktaAppPayload(appName, privateKey)
	if err != nil {
		return "", err
	}

	return makeRequest(requestURL, payload, token)
}

func createOktaAppPayload(appName string, privateKey *rsa.PrivateKey) (string, error) {
	key, err := jwk.New(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}

	keyJSON, err := jwkToJSON(key)
	if err != nil {
		return "", err
	}

	return formatOktaAppPayload(appName, keyJSON), nil
}

func formatOktaAppPayload(appName, keyJSON string) string {
	return fmt.Sprintf(`{
"name": "oidc_client",
"label": "%v",
"signOnMode": "OPENID_CONNECT",
"credentials": {
	"oauthClient": {
		"token_endpoint_auth_method": "private_key_jwt"
	}
},
"settings": {
	"oauthClient": {
		"redirect_uris": [
			"https://example.com"
		],
		"response_types": [
			"token"
		],
		"grant_types": [
			"client_credentials"
		],
		"application_type": "service",
		"jwks": {
			"keys": [ %v ]
		}
	}
}}`, appName, keyJSON)
}

func jwkToJSON(key jwk.Key) (string, error) {
	m := make(map[string]interface{})
	if err := key.PopulateMap(m); err != nil {
		return "", err
	}
	m["kid"] = "SIGNING_KEY"

	jsonbuf, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonbuf), nil
}
