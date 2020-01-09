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

func findZone(zones *hostingDeZoneConfigsFindResponse, domain string) *hostingDeZoneConfig {
	for _, zone := range zones.Response.Data {
		if zone.Name == domain {
			return &zone
		}
	}

	return nil
}

func createAuthRecord(rest *hostingDeRestAPI, domain string, validation string) error {
	domain = getTopLevelDomain(domain)

	zones, err := rest.zonesFind()
	if err != nil {
		return err
	}

	domainZone := findZone(zones, domain)
	if domainZone == nil {
		return fmt.Errorf("Could nor find '%s' domain", domain)
	}

	return nil
}
