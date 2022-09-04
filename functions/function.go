package functions

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sFTP_Transfer_Director/messages"
	"strings"
	"time"
)

var localFilePath string
var targetAddress string
var remoteFilePath string
var ActionFile string

var UploadAliases = []string{"upload", "u"}
var DownloadAliases = []string{"download", "d"}

func init() {
	flag.StringVar(&ActionFile, "a", "", "> Choose whether to download or upload our file")
	flag.StringVar(&localFilePath, "f", "", "> File path to transfer")
	flag.StringVar(&targetAddress, "t", os.Getenv("SFTP_TARGET"), "> Host Target Address")
	flag.StringVar(&remoteFilePath, "r", "", "> Remote file path to transfer/download to")
	flag.Parse()

	if localFilePath == "" || targetAddress == "" || remoteFilePath == "" || ActionFile == "" {
		log.Fatalln(messages.ErrorIncorrectSyntaxe)
	}
}

func TryToConnect(connectionState *bool) {
	point := 0
	tryToConnectMessage := "\rConnection in progress"

	for *connectionState == true {
		point += 1
		defineIntToString := strings.Repeat(".", point)

		fmt.Print(tryToConnectMessage + defineIntToString)
		time.Sleep(500 * time.Millisecond)
	}
}
