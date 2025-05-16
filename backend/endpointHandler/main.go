package main

import (
	"github.com/endpointhandler/repository"
	"github.com/endpointhandler/router"
	"log"
)

func main() {
	repository.InitDB()
	r := router.SetupRouter()
	log.Fatal(r.Run(":8000"))
}
