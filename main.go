package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	commandData, err := evalCommands();
	if err != nil {
		log.Fatal(err)
		return
	}

	var cfg config
	if err := readConfig(&cfg); err != nil {
		log.Fatal(err)
		return
	}

	rest := newhostingDeRestAPI(cfg.URL, cfg.AuthToken)

	var commandError error = nil;
	switch commandData.(type) {
	case *manualAuthCommand:
		manualAuthCmdData := commandData.(*manualAuthCommand)
		commandError = createAuthRecord(rest, manualAuthCmdData.domain, manualAuthCmdData.validation)

	case *manualCleanupCommand:
		manualCleanupCmdData := commandData.(*manualCleanupCommand)
		commandError = deleteAuthRecord(rest, manualCleanupCmdData.domain)

	case *daneUpdateCommand:
		daneUpdateCmdData := commandData.(*daneUpdateCommand)
		commandError = updateDaneRecord(rest, cfg, daneUpdateCmdData.domain, daneUpdateCmdData.certPath)

	default:
		commandError = fmt.Errorf("Unknown command data: %T", commandData)
	}

	if commandError != nil {
		log.Fatal(commandError)
		os.Exit(1)
	}
}
