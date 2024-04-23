package main

import (
	"multiple-servers/api"
	"flag"
	"fmt"
	"log"
	"os"
)

func main (){
	database := os.Getenv("DATABASE_URL")
	portPtr := flag.Int("port", 0, "Port number to pass static server to")
	flag.Parse()
	if database == "" {
		log.Fatalln("DATABASE_URL not set")
	}
    fmt.Println("Port:", *portPtr)
	api.DatabaseConnection(database, *portPtr)
}