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

type hostingDeKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type hostingDeTrace struct {
	Code          int                 `json:"code"`
	ContextObject string              `json:"contextObject"`
	ContextPath   string              `json:"contextPath"`
	Details       []hostingDeKeyValue `json:"details"`
	Text          string              `json:"text"`
	Value         string              `json:"value"`
}

type hostingDeMetadata struct {
	ClientTransactionID string `json:"clientTransactionId"`
	ServerTransactionID string `json:"serverTransactionId"`
}

type hostingDeResponse struct {
	Errors   []hostingDeTrace  `json:"errors"`
	Warnings []hostingDeTrace  `json:"warnings"`
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
	Type    string `json:"type"`
	Name	string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type hostingDeFullRecord struct {
	hostingDeRecord
	ID             string `json:"id"`
	AccountID      string `json:"accountId"`
	ZoneConfigID   string `json:"zoneConfigId"`
	AddDate        string `json:"addDate"`
	LastChangeDate string `json:"lastChangeDate"`
}

type hostingDeZoneConfigsFindRequest struct {
	AuthToken string `json:"authToken"`
}

type hostingDeZoneConfigsFindResponse struct {
	Response hostingDeZoneConfigsFindResponseData `json:"response"`
}

type hostingDeZoneConfigsFindResponseData struct {
	Data []hostingDeZoneConfig `json:"data"`
}

type hostingDeRecordsFindRequest struct {
	AuthToken string                            `json:"authToken"`
	Filter    hostingDeRecordsFindRequestFilter `json:"filter"`
	Limit     int                               `json:"limit"`
}

type hostingDeRecordsFindRequestFilter struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type hostingDeRecordsFindResponse struct {
	Response hostingDeRecordsFindResponseData `json:"response"`
}

type hostingDeRecordsFindResponseData struct {
	Data []hostingDeFullRecord `json:"data"`
}

type hostingDeZoneUpdateRequest struct {
	AuthToken       string                `json:"authToken"`
	ZoneConfig      hostingDeZoneConfig   `json:"zoneConfig"`
	RecordsToAdd    []hostingDeRecord     `json:"recordsToAdd"`
	RecordsToDelete []hostingDeFullRecord `json:"recordsToDelete"`
}

type hostingDeZoneUpdateResponseData struct {
	ZoneConfig hostingDeZoneConfig   `json:"zoneConfig"`
	Records    []hostingDeFullRecord `json:"records"`
}

type hostingDeZoneUpdateResponse struct {
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
	httpResponse, err := http.Post(url, "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return err
	}

	if httpResponse.StatusCode != 200 {
		return fmt.Errorf("HTTP error status: %d", httpResponse.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(httpResponse.Body)

	var restResponse hostingDeResponse
	err = json.Unmarshal(buf.Bytes(), &restResponse)
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

func (rest *hostingDeRestAPI) recordsFind(zoneConfigID string, limit int) (*hostingDeRecordsFindResponse, error) {
	request := &hostingDeRecordsFindRequest{rest.authToken, hostingDeRecordsFindRequestFilter{"zoneConfigId", zoneConfigID}, limit}
	response := new(hostingDeRecordsFindResponse)

	err := rest.call("recordsFind", request, response)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return response, nil
}

func (rest *hostingDeRestAPI) zoneUpdate(zoneConfig hostingDeZoneConfig, recordsToAdd []hostingDeRecord, recordsToDelete []hostingDeFullRecord) error {
	request := &hostingDeZoneUpdateRequest{rest.authToken, zoneConfig, recordsToAdd, recordsToDelete}
	response := new(hostingDeZoneUpdateResponse);

	err := rest.call("zoneUpdate", request, response)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
