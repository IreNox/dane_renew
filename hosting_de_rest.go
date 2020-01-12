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
	ClientTransactionID string `json:"clientTransactionId"`
	ServerTransactionID string `json:"serverTransactionId"`
}

type hostingDeZoneConfigSoaValues struct {
	Refresh     int `json:"refresh"`
	Retry       int `json:"retry"`
	Expire      int `json:"expire"`
	TTL         int `json:"ttl"`
	NegativeTTL int `json:"negativeTtl"`
}

type hostingDeZoneConfig struct {
	ID                    string                       `json:"id"`
	AccountID             string                       `json:"accountId"`
	DNSSecMode            string                       `json:"dnsSecMode"`
	EmailAddress          string                       `json:"emailAddress"`
	AddDate               string                       `json:"addDate"`
	LastChangeDate        string                       `json:"lastChangeDate"`
	MasterIP              string                       `json:"masterIp"`
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

type hostingDeRecord struct {
	//ID      string `json:"id"`
	Type    string `json:"type"`
	Name	string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type hostingDeZoneUpdateRequest struct {
	AuthToken       string              `json:"authToken"`
	ZoneConfig      hostingDeZoneConfig `json:"zoneConfig"`
	RecordsToAdd    []hostingDeRecord   `json:"recordsToAdd"`
	RecordsToDelete []hostingDeRecord   `json:"recordsToDelete"`
}

type hostingDeZoneUpdateResponseData struct {
	ZoneConfig hostingDeZoneConfig `json:"zoneConfig"`
	Records    []hostingDeRecord     `json:"records"`
}

type hostingDeZoneUpdateResponse struct {
	Errors   []string                        `json:"errors"`
	Warnings []string                        `json:"warnings"`
	Status   string                          `json:"status"`
	Metadata hostingDeMetadata               `json:"metadata"`
	Response hostingDeZoneUpdateResponseData `json:"response"`
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

func (rest *hostingDeRestAPI) call(function string, request interface{}, response interface{}) error {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := rest.url + function
	fmt.Printf("Request URL: %s, Body: %s\n", url, string(requestBytes))

	httpResponse, err := http.Post(url, "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return err
	}

	if httpResponse.StatusCode != 200 {
		return fmt.Errorf("HTTP error status: %d", httpResponse.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(httpResponse.Body)

	ioutil.WriteFile("D:\\test.json", buf.Bytes(), 0777)

	err = json.Unmarshal(buf.Bytes(), response)
	if err != nil {
		return err
	}

	return nil
}

func (rest *hostingDeRestAPI) zonesFind() (*hostingDeZoneConfigsFindResponse, error) {
	request := &hostingDeZoneConfigsFindRequest{rest.authToken}
	response := new(hostingDeZoneConfigsFindResponse)

	err := rest.call("zoneConfigsFind", request, response)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return response, nil
}

func (rest *hostingDeRestAPI) zoneUpdate(zoneConfig hostingDeZoneConfig, recordsToAdd []hostingDeRecord, recordsToDelete []hostingDeRecord) error {
	request := &hostingDeZoneUpdateRequest{rest.authToken, zoneConfig, recordsToAdd, recordsToDelete}
	response := new(hostingDeZoneUpdateResponse);

	err := rest.call("zoneUpdate", request, response)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
