package main

import "strings"

import "fmt"

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

func createAuthRecord(rest *hostingDeRestAPI, domain string, validation string) error {
	domainZoneConfig, err := findZone(rest, domain)
	if err != nil {
		return err
	}

	recordsToAdd := []hostingDeRecord{hostingDeRecord{"authChallenge", "TXT", "_acme-challenge." + domain, validation, 120}}
	recordsToDelete := []hostingDeRecord{}
	err = rest.zoneUpdate(*domainZoneConfig, recordsToAdd, recordsToDelete)
	if err != nil {
		return err
	}

	return nil
}

func deleteAuthRecord(rest* hostingDeRestAPI, domain string, validation string) error {
	domainZoneConfig, err := findZone(rest, domain)
	if err != nil {
		return err
	}

	recordsToAdd := []hostingDeRecord{}
	recordsToDelete := []hostingDeRecord{} //hostingDeRecord{"TXT", "_acme-challenge." + domain, validation, 120}}
	err = rest.zoneUpdate(*domainZoneConfig, recordsToAdd, recordsToDelete)
	if err != nil {
		return err
	}

	return nil
}
