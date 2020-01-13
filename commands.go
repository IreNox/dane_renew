package main

import (
	"flag"
	"fmt"
	"os"
)

type commandsArgument struct {
	name      string
	envName   string
	argsValue *string
	target    *string
}

type commandsCommand struct {
	name      string
	flagSet   *flag.FlagSet
	arguments []commandsArgument
	target    interface{}
}

type commandsContext struct {
	commands []*commandsCommand
}

type manualAuthCommand struct {
	domain     string
	validation string
}

type manualCleanupCommand struct {
	domain string
}

func (context *commandsContext) addCommand(name string, target interface{}) *commandsCommand {
	command := &commandsCommand{name: name, target: target}
	command.flagSet = flag.NewFlagSet(name, flag.ExitOnError)
	context.commands = append(context.commands, command)
	return command
}

func (context *commandsContext) parse() (interface{}, error) {
	commandName := os.Args[1]
	for _, command := range context.commands {
		if command.name != commandName {
			continue
		}

		err := command.flagSet.Parse(os.Args[2:])
		if err != nil {
			return nil, err
		}

		for _, argument := range command.arguments {
			if *argument.argsValue == "" {
				envValue, exists := os.LookupEnv(argument.envName)
				if !exists {
					return nil, fmt.Errorf("No argument '-%s' and no '%s' environment variable found", argument.name, argument.envName)
				}
				*argument.target = envValue
			} else {
				*argument.target = *argument.argsValue
			}
		}

		return command.target, nil
	}

	return nil, fmt.Errorf("Unknown command: %s", commandName)
}

func (command *commandsCommand) addArgument(name string, description string, envName string, target *string) {
	argument := commandsArgument{name: name, envName: envName, target: target}
	argument.argsValue = command.flagSet.String(name, "", description)
	command.arguments = append(command.arguments, argument)
}

func evalCommands() (interface{}, error) {
    if len(os.Args) < 2 {
		return nil, fmt.Errorf("Not enought arguments %d needs 2", len(os.Args))
	}

	var context commandsContext

	manualAuthCmdData := new(manualAuthCommand)
	manualAuthCmd := context.addCommand("manual-auth", manualAuthCmdData)
	manualAuthCmd.addArgument("domain", "domain name", "CERTBOT_DOMAIN", &manualAuthCmdData.domain)
	manualAuthCmd.addArgument("validation", "validation string", "CERTBOT_VALIDATION", &manualAuthCmdData.validation);

	manualCleanupCmdData := new(manualCleanupCommand)
	manualCleanupCmd := context.addCommand("manual-cleanup", manualCleanupCmdData)
	manualCleanupCmd.addArgument("domain", "domain name", "CERTBOT_DOMAIN", &manualCleanupCmdData.domain)

	return context.parse()
}