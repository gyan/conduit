package main

import (
	"flag"
	"fmt"
	"os"
	"strings"


	_ "github.com/joho/godotenv/autoload"
	pm "github.com/arpit32/conduit/pendulum"
)

const (
	workflow = "workflow"
	activity = "activity"
)

func main() {
	var service string
	var configFilePath string
	var serverPort string
	var tasklist string
	var logger string
	flag.StringVar(&service, "service", "", "Name of the service to start")
	flag.StringVar(&configFilePath, "config", "./config", "path to the configuration file")
	flag.StringVar(&serverPort, "port", "4000", "port on which server runs")
	flag.StringVar(&tasklist, "tasklist", "", "Name of the tasklist")
	flag.StringVar(&logger, "v", "0", "Logger enable/disable")
	flag.Parse()

	if os.Getenv("SERVICE") != "" {
		service = os.Getenv("SERVICE")
	}

	if os.Getenv("TASKLIST") != "" {
		tasklist = os.Getenv("TASKLIST")
	}

	if os.Getenv("VERBOSE") != "" {
		logger = os.Getenv("VERBOSE")
	}

	fmt.Println(service, tasklist)

	if service == "app" || os.Getenv("SERVICE") == "app" {
		// application := server.New(configFilePath)
		// application.Init()
		// application.Start(serverPort)
	} else {
		if strings.ToLower(tasklist) == strings.ToLower(pm.TaskList) {
			
			if service == workflow || os.Getenv("SERVICE") == workflow {
				fmt.Println("Workflow worker starting")
				worker := pm.New(configFilePath)
				worker.Init(tasklist, logger, service)
				worker.Start()
			} else if service == activity || os.Getenv("SERVICE") == activity {
				fmt.Println("Activity worker starting")
				worker := pm.New(configFilePath)
				worker.Init(tasklist, logger, service)
				worker.Start()
			}
		}
	}
}
