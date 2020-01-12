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

type hostingDeErrorDetail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type hostingDeError struct {
	Code          int                    `json:"code"`
	ContextObject string                 `json:"contextObject"`
	ContextPath   string                 `json:"contextPath"`
	Details       []hostingDeErrorDetail `json:"details"`
	Text          string                 `json:"text"`
	Value         string                 `json:"value"`
}

type hostingDeMetadata struct {
	ClientTransactionID string `json:"clientTransactionId"`
	ServerTransactionID string `json:"serverTransactionId"`
}

type hostingDeResponse struct {
	Errors   []hostingDeError  `json:"errors"`
	Warnings []string          `json:"warnings"`
	Status   string            `json:"status"`
	Metadata hostingDeMetadata `json:"metadata"`
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

type hostingDeRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name	string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type hostingDeZoneConfigsFindRequest struct {
	AuthToken string `json:"authToken"`
}

type hostingDeZoneConfigsFindResponseData struct {
	Data []hostingDeZoneConfig `json:"data"`
}

type hostingDeZoneConfigsFindResponse struct {
	Response hostingDeZoneConfigsFindResponseData `json:"response"`
}

type hostingDeRecordsFindRequest struct {
	AuthToken string                            `json:"authToken"`
	Filter    hostingDeRecordsFindRequestFilter `json:"filter"`
}

type hostingDeRecordsFindRequestFilter struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type hostingDeRecordsFindResponse struct {
	Response hostingDeZoneConfigsFindResponseData `json:"response"`
}

type hostingDeRecordsFindResponseData struct {
	Data []hostingDeZoneConfig `json:"data"`
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

	var restResponse hostingDeResponse
	err = json.Unmarshal(buf.Bytes(), restResponse)
	if err != nil {
		return err
	}

	if len(restResponse.Errors) > 0 {
		return fmt.Errorf(restResponse.Errors[0].Text)
	}

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

func (rest *hostingDeRestAPI) recordsFind(zoneConfigId string) (*hostingDeRecordsFindResponse, error) {	
	request := &hostingDeRecordsFindRequest{rest.authToken, hostingDeRecordsFindRequestFilter{"zoneConfigId", zoneConfigId}}
	response := new(hostingDeRecordsFindResponse)

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
