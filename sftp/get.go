package sftp

import (
	"log"
	"os"
	"path"
	"sync"
)

func Get(user, password, key, host, port, dst string, remoteFiles []string) bool {
	sftpClient, err := Connect(user, password, key, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	var wg sync.WaitGroup

	for _, file := range remoteFiles {
		wg.Add(1)

		go func(file string) {
			defer wg.Add(-1)

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
		}(file)

/*
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
*/
	}
	wg.Wait()

	return true
}
