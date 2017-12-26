package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func initLogging(isDebugEnabled bool) {
	log.SetOutput(os.Stdout)

	if isDebugEnabled {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	// flags
	configFile := kingpin.Flag("config.file", "Azure configuration file").Default("sample-azure.yml").String()
	isDebugEnabled := kingpin.Flag("debug", "Debug flag").Default("false").Bool()

	// get command
	getCommand := kingpin.Command("get", "Perform GET <url> against Azure Resource Manager API")
	getCommandUrl := getCommand.Arg("url", "The <url>").Required().String()

	// summary command
	summaryCommand := kingpin.Command("resources", "Print out the Azure resources that exist on this subscription")
	summaryCommandMaxContinuation := summaryCommand.Arg("maxcontinuation", "The max number of continuations to follow.  Default to 10.").Default("10").Int()

	command := kingpin.Parse()

	// initialize logging after parsing flags
	initLogging(*isDebugEnabled)

	config := &Config{}
	err := config.loadConfig(*configFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	environment := getCurrentEnvironment(config.Credentials.Environment)
	processor := NewCommandProcessor(config, environment)

	// process commands
	switch command {
	case "get":
		processor.processGetCommand(*getCommandUrl)
		break
	case "resources":
		processor.processSummarizeCommand(*summaryCommandMaxContinuation)
		break
	default:
		log.Errorf("Unknown command: %s\n", command)
		break
	}
}