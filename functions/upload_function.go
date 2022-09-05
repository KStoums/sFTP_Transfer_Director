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

func UploadSFTP() {
	file, err := os.Open(localFilePath)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			log.Fatalln(messages.ErrorNoPermissionToOpenFile)
		}
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalln(messages.ErrorFileNotExist)
		}
		if errors.Is(err, os.ErrInvalid) {
			log.Fatalln(messages.ErrorInvalidPath)
		}
		log.Fatalln(err)
	}
	defer file.Close()

	fmt.Print(messages.EnterUsername)
	var username string
	fmt.Scanln(&username)

	fmt.Print(messages.EnterPassword)
	passwordByte, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	password := string(passwordByte)
	fmt.Println("")

	tryToConnectState := true
	go TryToConnect(&tryToConnectState)

	sshclient, err := ssh.Dial("tcp", targetAddress, &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	})

	tryToConnectState = false
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

	if strings.ContainsAny(localFilePath, "\\") {
		strings.ReplaceAll(localFilePath, "\\", "/")
	}

	splitted := strings.Split(localFilePath, "/")
	localFileName := splitted[len(splitted)-1]

	fileInfo, err := sftpClient.Stat(remoteFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = sftpClient.MkdirAll(remoteFilePath)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln(err)
		}
	}

	localFileInfo, _ := file.Stat()

	if fileInfo.IsDir() {
		if !strings.HasSuffix(remoteFilePath, "/") && !strings.HasSuffix(remoteFilePath, "\\") {
			remoteFilePath += "/"
		}
		remoteFilePath += localFileName
	}

	remoteFile, err := sftpClient.OpenFile(remoteFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
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

	bar := progressbar.NewOptions64(localFileInfo.Size(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]Uploading %s...[reset]", localFileName)),
	)

	if _, err = io.Copy(io.MultiWriter(remoteFile, bar), file); err != nil {
		log.Fatalln(err)
	}
}
