package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type appResponse struct {
	Credentials credentials `json:"credentials"`
}
type credentials struct {
	OAuthClient oauthClient `json:"oauthClient"`
}
type oauthClient struct {
	ClientID string `json:"client_id"`
}

func makeRequest(requestURL, payload, token string) (string, error) {
	req, err := http.NewRequest("POST", requestURL, bytes.NewReader([]byte(payload)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	tokenHeader := fmt.Sprintf("SSWS %s", token)
	req.Header.Set("Authorization", tokenHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cache-control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	appResponse, err := processResponse(resp)
	if err != nil {
		return "", err
	}

	return appResponse.Credentials.OAuthClient.ClientID, nil
}

func processResponse(resp *http.Response) (*appResponse, error) {
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("App create failed with status code: %v\n\n%v",
			resp.StatusCode, string(responseData))
	}

	appResponse := &appResponse{}
	err = json.Unmarshal(responseData, appResponse)
	if err != nil {
		return nil, err
	}

	return appResponse, nil
}
