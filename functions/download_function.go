package functions

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"sFTP_Transfer_Director/messages"
	"strings"
	"syscall"
	"time"
)

func DownloadSFTP() {
	fmt.Print(messages.EnterUsername)
	var username string
	fmt.Scanln(&username)

	fmt.Print(messages.EnterPassword)
	passwordByte, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	password := string(passwordByte)
	fmt.Println("")

	TryToConnectState := true
	go TryToConnect(&TryToConnectState)

	sshclient, err := ssh.Dial("tcp", targetAddress, &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	})

	TryToConnectState = false
	fmt.Println()
	if err != nil {
		log.Fatalln(err)
	}
	defer sshclient.Close()

	sftpClient, err := sftp.NewClient(sshclient)
	if err != nil {
		log.Fatalln(err)
	}
	defer sftpClient.Close()

	strings.ReplaceAll(localFilePath, "\\", "/")

	splitted := strings.Split(remoteFilePath, "/")
	remoteFileName := splitted[len(splitted)-1]

	remoteFile, err := sftpClient.OpenFile(remoteFilePath, os.O_RDONLY)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			log.Fatalln(messages.ErrorNoPermissionToOpenFile)
		} else if errors.Is(err, os.ErrInvalid) {
			log.Fatalln(messages.ErrorInvalidPath)
		} else {
			log.Fatalln(err)
		}
	}
	defer remoteFile.Close()
	localFileInfo, err := os.Stat(localFilePath)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			log.Fatalln(messages.ErrorNoPermissionToOpenFile)
		}
		if errors.Is(err, os.ErrInvalid) {
			log.Fatalln(messages.ErrorInvalidPath)
		}
		log.Fatalln(err)
	}

	if localFileInfo.IsDir() {
		if !strings.HasSuffix(localFilePath, "/") && !strings.HasSuffix(localFilePath, "\\") {
			localFilePath += "/"
		}
		localFilePath += remoteFileName
	}

	localFile, err := os.OpenFile(localFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			log.Fatalln(messages.ErrorNoPermissionToOpenFile)
		}
		if errors.Is(err, os.ErrInvalid) {
			log.Fatalln(messages.ErrorInvalidPath)
		}
		log.Fatalln(err)
	}

	remoteFileInfo, _ := remoteFile.Stat()
	bar := progressbar.NewOptions64(remoteFileInfo.Size(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]Downloading %s...[reset]", remoteFileInfo.Name())),
	)

	if _, err = io.Copy(io.MultiWriter(localFile, bar), remoteFile); err != nil {
		log.Fatalln(err)
	}
}
