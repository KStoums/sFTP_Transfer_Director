package main

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
	"strings"
	"syscall"
	"time"
)

func main() {
	file, err := os.Open(localFilePath)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			log.Fatalln("Error: No permission to open file.")
		}
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalln("Error: File does not exist.")
		}
		if errors.Is(err, os.ErrInvalid) {
			log.Fatalln("Error: Invalid path.")
		}
		log.Fatalln(err)
	}
	defer file.Close()

	fmt.Print("Enter server username > ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("Enter server password > ")
	passwordByte, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	password := string(passwordByte)
	fmt.Println("")

	tryToConnectState := true
	go tryToConnect(&tryToConnectState)

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

	/*path := "/home/test/Kbot"
	  splitted := strings.Split(path, "/")
	  fmt.Println(splitted[len(splitted)-1])
	*/

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

	fileInfo, _ = sftpClient.Stat(remoteFilePath)
	localFileInfo, _ := file.Stat()

	if fileInfo.IsDir() {
		if !strings.HasSuffix(remoteFilePath, "/") && !strings.HasSuffix(remoteFilePath, "\\") {
			remoteFilePath += "/"
		}
		remoteFilePath += localFileInfo.Name()
	}

	remoteFile, err := sftpClient.OpenFile(remoteFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			log.Fatalln("Error: No permission to open file.")
		} else if errors.Is(err, os.ErrInvalid) {
			log.Fatalln("Error: Invalid path.")
		} else {
			log.Fatalln(err)
		}
	}
	defer remoteFile.Close()

	bar := progressbar.DefaultBytes(
		localFileInfo.Size(),
		"downloading",
	)
	io.Copy(io.MultiWriter(remoteFile, bar), file)
}
