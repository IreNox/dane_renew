package main

import (
	"fmt"
	"log"
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
	fmt.Printf("Config: %s\n", cfg)

	rest := newhostingDeRestAPI(cfg.URL, cfg.AuthToken)

	var commandError error = nil;
	switch commandData.(type) {
	case *manualAuthCommand:
		manualAuthCmdData := commandData.(*manualAuthCommand)
		commandError = createAuthRecord(rest, manualAuthCmdData.domain, manualAuthCmdData.validation)

	case *manualCleanupCommand:
		manualCleanupCmdData := commandData.(*manualCleanupCommand)
		commandError = deleteAuthRecord(rest, manualCleanupCmdData.domain)

	default:
		commandError = fmt.Errorf("Unknown command data: %T", commandData)
	}

	if commandError != nil {
		log.Fatal(commandError)
		return
	}
}
