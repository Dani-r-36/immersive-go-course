package main

import (
	"multiple-servers/static"
	"flag"
	"fmt"
	"path/filepath"
	"log"
)

func main (){
	pathPtr := flag.String("path", "", "path for website assets")
	portPtr := flag.Int("port", 0, "Port number to pass static server to")
	flag.Parse()
	fmt.Println("Path:", *pathPtr)
    fmt.Println("Port:", *portPtr)
	absPath, err := filepath.Abs(*pathPtr)
	if err != nil {
		log.Fatalln(err)
	}
	static.Run(static.Config{absPath, *portPtr})
}