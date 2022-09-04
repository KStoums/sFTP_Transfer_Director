package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var localFilePath string
var targetAddress string
var remoteFilePath string

func init() {
	flag.StringVar(&localFilePath, "f", "", "> File path to transfer")
	flag.StringVar(&targetAddress, "t", os.Getenv("SFTP_TARGET"), "> Host Target Address")
	flag.StringVar(&remoteFilePath, "r", "", "> Remote file path to transfer to")
	flag.Parse()

	if localFilePath == "" || targetAddress == "" || remoteFilePath == "" {
		log.Fatalln("Error: Incorrect syntax.")
	}
}

func tryToConnect(connectionState *bool) {
	point := 0
	tryToConnectMessage := "\rConnection in progress"

	for *connectionState == true {
		point += 1
		defineIntToString := strings.Repeat(".", point)

		fmt.Print(tryToConnectMessage + defineIntToString)
		time.Sleep(500 * time.Millisecond)
	}
}
