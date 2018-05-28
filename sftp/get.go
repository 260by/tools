package sftp

import (
	// "fmt"
	"log"
	"os"
	"path"
)

func Get(user, password, key, host, port, dst string, remoteFiles []string) bool {
	sftpClient, err := Connect(user, password, key, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()
	for _, file := range remoteFiles {
		var remoteFilePath = file
		var localDir = dst

		srcFile, err := sftpClient.Open(remoteFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer srcFile.Close()

		var localFileName = path.Base(remoteFilePath)
		dstFile, err := os.Create(path.Join(localDir, localFileName))
		if err != nil {
			log.Fatal(err)
		}
		defer dstFile.Close()

		if _, err = srcFile.WriteTo(dstFile); err != nil {
			log.Fatal(err)
		}
	}

	return true
}
