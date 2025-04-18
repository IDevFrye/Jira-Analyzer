package main

import (
	"log"
	"os"

	"github.com/jiraconnector/cmd/app"
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
	a, err := app.NewApp(cfg)
	if err != nil {
		log.Println("error create app")
		panic(err)
	}
	log.Println("created app")

	//start app
	if err := a.Run(); err != nil {
		log.Println("error run app")
		panic(err)
	}
	defer a.Close()
}
