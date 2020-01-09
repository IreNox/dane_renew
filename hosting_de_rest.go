package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type hostingDeRestAPI struct {
	url       string
	authToken string
}

type hostingDeZoneConfigsFindRequest struct {
	AuthToken string `json:"authToken"`
}

type hostingDeMetadata struct {
	ClientTransactionId string `json:"clientTransactionId"`
	ServerTransactionId string `json:"serverTransactionId"`
}

type hostingDeZoneConfigSoaValues struct {
	Refresh     int `json:"refresh"`
	Retry       int `json:"retry"`
	Expire      int `json:"expire"`
	Ttl         int `json:"ttl"`
	NegativeTtl int `json:"negativeTtl"`
}

type hostingDeZoneConfig struct {
	Id                    string                       `json:"id"`
	AccountId             string                       `json:"accountId"`
	DnsSecMode            string                       `json:"dnsSecMode"`
	EmailAddress          string                       `json:"emailAddress"`
	AddDate               string                       `json:"addDate"`
	LastChangeDate        string                       `json:"lastChangeDate"`
	MasterIp              string                       `json:"masterIp"`
	Name                  string                       `json:"name"`
	NameUnicode           string                       `json:"nameUnicode"`
	SoaValues             hostingDeZoneConfigSoaValues `json:"soaValues"`
	Status                string                       `json:"status"`
	Type                  string                       `json:"type"`
	ZoneTransferWhitelist []string                     `json:"zoneTransferWhitelist"`
}

type hostingDeZoneConfigsFindResponseData struct {
	Data []hostingDeZoneConfig `json:"data"`
}

type hostingDeZoneConfigsFindResponse struct {
	Errors   []string                             `json:"errors"`
	Warnings []string                             `json:"warnings"`
	Status   string                               `json:"status"`
	Metadata hostingDeMetadata                    `json:"metadata"`
	Response hostingDeZoneConfigsFindResponseData `json:"response"`
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

func (rest *hostingDeRestAPI) call(function string, request interface{}) ([]byte, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := rest.url + function
	fmt.Printf("Request URL: %s, Body: %s\n", url, string(requestBytes))

	response, err := http.Post(url, "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error status: %d", response.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	ioutil.WriteFile("D:\\test.json", buf.Bytes(), 0777)
	return buf.Bytes(), nil
}

func (rest *hostingDeRestAPI) zonesFind() (*hostingDeZoneConfigsFindResponse, error) {
	request := &hostingDeZoneConfigsFindRequest{rest.authToken}

	responseBytes, err := rest.call("zoneConfigsFind", request)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	response := new(hostingDeZoneConfigsFindResponse)
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return response, nil
}
