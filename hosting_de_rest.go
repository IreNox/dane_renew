package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type hostingDeRestAPI struct {
	url       string
	authToken string
}

type hostingDeZonesFindRequest struct {
	AuthToken string `json:"authToken"`
}

func newhostingDeRestAPI(url string, authToken string) *hostingDeRestAPI {
	rest := new(hostingDeRestAPI)
	rest.url = url
	rest.authToken = authToken

	if !strings.HasSuffix(rest.url, "/") {
		rest.url += "/"
	}

	return rest
}

func (rest *hostingDeRestAPI) call(function string, request interface{}) (string, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	url := rest.url + function
	fmt.Printf("Request URL: %s, Body: %s\n", url, string(requestBytes))

	response, err := http.Post(url, "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("HTTP error status: %d", response.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return buf.String(), nil
}

func (rest *hostingDeRestAPI) zonesFind() (string, error) {
	request := &hostingDeZonesFindRequest{rest.authToken}

	response, err := rest.call("zoneConfigsFind", request)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return response, nil
}
