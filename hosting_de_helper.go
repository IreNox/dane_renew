package main

import (
	"strings"
	"fmt"
	"time"
)

func getTopLevelDomain(domain string) string {
	for strings.Count(domain, ".") > 1 {
		index := strings.Index(domain, ".")
		domain = domain[index+1:]
	}

	return domain
}

func requestAndFindZone(rest *hostingDeRestAPI, domain string) (*hostingDeZoneConfig, error) {
	topLevelDomain := getTopLevelDomain(domain)

	zones, err := rest.zonesFind()
	if err != nil {
		return nil, err
	}

	for _, zone := range zones.Response.Data {
		if zone.Name == topLevelDomain {
			return &zone, nil
		}
	}

	return nil, fmt.Errorf("Could not find '%s' domain", topLevelDomain)
}

func requestAndFindSafeZone(rest *hostingDeRestAPI, domain string) (*hostingDeZoneConfig, error) {
	zone, err := requestAndFindZone(rest, domain)
	if err != nil {
		return nil, err
	}

	var blockedTimer int = 0
	for zone.Status == "blocked" {
		fmt.Printf("Zone '%s' is blocked wait for 10 seconds.\n", zone.Name)

		blockedTimer++
		if blockedTimer == 10 {
			return nil, fmt.Errorf("Zone '%s' is blocked", zone.Name)
		}

		time.Sleep(10 * time.Second)

		zone, err = requestAndFindZone(rest, domain)
		if err != nil {
			return nil, err
		}
	}

	return zone, nil
}

func requestRecords(rest *hostingDeRestAPI, domain string) (*hostingDeZoneConfig, *hostingDeRecordsFindResponse, error) {
	zoneConfig, err := requestAndFindSafeZone(rest, domain)
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
