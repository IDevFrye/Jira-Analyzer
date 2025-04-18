package main

import (
	"log"
	"os"

	config "github.com/jiraconnector/internal/configReader"
)

func main() {
	//Open config
	cfgPath := "../../configs/config.yml" // dev
	configFile, err := os.Open(cfgPath)
	if err != nil {
		log.Println("error open config")
		panic(err)
	}
	log.Println("open config")

	//load data from config
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Println("error load config")
		panic(err)
	}
	log.Println("load config")

	//setting logger
	log.Printf("set logger") //TEMP

	//create connector app

	//start app
}
