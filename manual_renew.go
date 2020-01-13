package main

func createAuthRecord(rest *hostingDeRestAPI, domain string, validation string) error {
	domainZoneConfig, err := requestAndFindSafeZone(rest, domain)
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
	domainZoneConfig, record, err := requestAndFindRecord(rest, "_acme-challenge." + domain, "TXT")
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
