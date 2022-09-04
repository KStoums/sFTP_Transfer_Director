package main

import (
	"github.com/thoas/go-funk"
	"log"
	"sFTP_Transfer_Director/functions"
	"sFTP_Transfer_Director/messages"
)

func main() {
	if funk.ContainsString(functions.UploadAliases, functions.ActionFile) {
		functions.UploadSFTP()
		return
	}

	if funk.ContainsString(functions.DownloadAliases, functions.ActionFile) {
		functions.DownloadSFTP()
		return
	}

	log.Fatalln(messages.ErrorIncorrectSyntaxe)
}
