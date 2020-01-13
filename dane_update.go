package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"fmt"
	"strconv"
	"os"
)

func calculateCertHash(certPath string) (string, error) {
	certFile, err := os.Open(certPath)
	if err != nil {
		return "", err
	}
	defer certFile.Close()

	certBytes, err := ioutil.ReadAll(certFile)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(certBytes)
	if block == nil {
		return "", fmt.Errorf("Failed to de code PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", err
	}

	publicKeyHash := sha256.Sum256(publicKeyBytes)
	publicKeyHashHex := hex.EncodeToString(publicKeyHash[:])

	return publicKeyHashHex, nil
}

func updateDaneRecord(rest* hostingDeRestAPI, cfg config, domain string, certPath string) error {
	publicKeyHashHex, err := calculateCertHash(certPath)
	if err != nil {
		return err
	}

	domainZoneConfig, records, err := requestRecords(rest, domain)
	if err != nil {
		return err
	}

	recordsToAdd := []hostingDeRecord{}
	recordsToDelete := []hostingDeFullRecord{}
	for _, cfgDomain := range cfg.Domains {
		if cfgDomain.Name != domain {
			continue
		}

		daneDomain := "_" + strconv.Itoa(cfgDomain.Port) + "._" + cfgDomain.Protocol + "." + domain
		daneValue := "3 1 1 " + publicKeyHashHex

		record, err := findRecord(records, daneDomain, "TLSA")
		if err != nil {
			return err
		}

		if record.Content == daneValue {
			fmt.Printf("Record '%s' is already up to date.\n", daneDomain )
			continue
		}

		recordsToAdd = append(recordsToAdd, hostingDeRecord{"TLSA", daneDomain, daneValue, 86400});
		recordsToDelete = append(recordsToDelete, *record)
	}

	if(len(recordsToAdd) == 0 && len(recordsToDelete) == 0) {
		return nil
	}

	return rest.zoneUpdate(*domainZoneConfig, recordsToAdd, recordsToDelete)
}