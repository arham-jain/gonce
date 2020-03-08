package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gonce/config"
	"github.com/gonce/services"
	"github.com/gonce/utils"
)

var d Dependencies

// Dependencies struct used to aggregate multiple dependencies
type Dependencies struct {
	ConditionConfig config.ConditionConfig
	HTTPClient      utils.HTTPClient
	PerformService  services.PerformService
}

func init() {
	initConfig()
	initHTTPClient()
}

func main() {
	log.Println(d.ConditionConfig)
	d.PerformService.Perform(d.ConditionConfig, d.HTTPClient)
	byteResult, _ := json.Marshal(d.PerformService.Metrics)
	log.Println(string(byteResult))
}

func initConfig() {
	appConfigLocation := os.Args[1]
	configFile, err := ioutil.ReadFile(appConfigLocation)
	if err != nil {
		log.Fatalf("Error reading condition config file : %s", appConfigLocation)
	}
	err = json.Unmarshal(configFile, &d.ConditionConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling condition config : %s", err.Error())
	}
}

func initHTTPClient() {
	d.HTTPClient.Initialise()
}
