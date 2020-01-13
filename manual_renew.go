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

func findRecord(rest *hostingDeRestAPI, domain string, recordType string) (*hostingDeZoneConfig, *hostingDeFullRecord, error) {
	domainZoneConfig, err := findZone(rest, domain)
	if err != nil {
		return nil, nil, err
	}

	records, err := rest.recordsFind(domainZoneConfig.ID)
	if err != nil {
		return nil, nil, err
	}

	for _, record := range records.Response.Data {
		if record.Name == domain && record.Type == recordType {
			return domainZoneConfig, &record, nil
		}
	}

	return nil, nil, fmt.Errorf("Could not find record '%s' of type '%s'", domain, recordType)
}

func createAuthRecord(rest *hostingDeRestAPI, domain string, validation string) error {
	domainZoneConfig, err := findZone(rest, domain)
	if err != nil {
		return err
	}

	recordsToAdd := []hostingDeRecord{hostingDeRecord{"TXT", "_acme-challenge." + domain, validation, 120}}
	recordsToDelete := []hostingDeFullRecord{}
	err = rest.zoneUpdate(*domainZoneConfig, recordsToAdd, recordsToDelete)
	if err != nil {
		return err
	}

	return nil
}

func deleteAuthRecord(rest* hostingDeRestAPI, domain string) error {
	domainZoneConfig, record, err := findRecord(rest, "_acme-challenge." + domain, "TXT")
	if err != nil {
		return err
	}

	recordsToAdd := []hostingDeRecord{}
	recordsToDelete := []hostingDeFullRecord{*record}
	err = rest.zoneUpdate(*domainZoneConfig, recordsToAdd, recordsToDelete)
	if err != nil {
		return err
	}

	return nil
}
