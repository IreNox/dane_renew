package main

import (
	"strings"
	"fmt"
)

func getTopLevelDomain(domain string) string {
	for strings.Count(domain, ".") > 1 {
		index := strings.Index(domain, ".")
		domain = domain[index+1:]
	}

	return domain
}

func findZone(rest *hostingDeRestAPI, domain string) (*hostingDeZoneConfig, error) {
	topLevelDomain := getTopLevelDomain(domain)

	zones, err := rest.zonesFind()
	if err != nil {
		return nil,err
	}

	for _, zone := range zones.Response.Data {
		if zone.Name == topLevelDomain {
			return &zone, nil
		}
	}

	return nil, fmt.Errorf("Could not find '%s' domain", topLevelDomain)
}

func requestRecords(rest *hostingDeRestAPI, domain string) (*hostingDeZoneConfig, *hostingDeRecordsFindResponse, error) {
	zoneConfig, err := findZone(rest, domain)
	if err != nil {
		return nil, nil, err
	}

	records, err := rest.recordsFind(zoneConfig.ID, 100)
	if err != nil {
		return nil, nil, err
	}

	return zoneConfig, records, nil
}

func findRecord(records *hostingDeRecordsFindResponse, domain string, recordType string) (*hostingDeFullRecord, error) {
	for _, record := range records.Response.Data {
		if record.Name == domain && record.Type == recordType {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("Could not find record '%s' of type '%s'", domain, recordType)
}

func requestAndFindRecord(rest *hostingDeRestAPI, domain string, recordType string) (*hostingDeZoneConfig, *hostingDeFullRecord, error) {
	zoneConfig, records, err := requestRecords(rest, domain)
	if err != nil {
		return nil, nil, err
	}

	record, err := findRecord(records, domain, recordType)
	if err != nil {
		return nil, nil, err
	}

	return zoneConfig, record, nil
}
