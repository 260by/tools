package sftp

import (
	"os"
	"path"
	"sync"
	"github.com/pkg/sftp"
)

func Get(sftpClient *sftp.Client, src, dst string, remoteFiles []string) (result bool, err error) {
	// sftpClient, err := Connect(user, password, key, host, port)
	// if err != nil {
	// 	log.Fatal(err)
	// }


	var wg sync.WaitGroup

	for _, file := range remoteFiles {
		wg.Add(1)

		go func(file string) {
			defer wg.Add(-1)

			var remoteFilePath = path.Join(src, file)
			println(remoteFilePath)
			var localDir = dst

			srcFile, err := sftpClient.Open(remoteFilePath)
			if err != nil {
				return
			}
			defer srcFile.Close()

			var localFileName = path.Base(remoteFilePath)
			dstFile, err := os.Create(path.Join(localDir, localFileName))
			if err != nil {
				return
			}
			defer dstFile.Close()

			if _, err = srcFile.WriteTo(dstFile); err != nil {
				return
			}
		}(file)

		// defer sftpClient.Close()
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

	defer sftpClient.Close()
	return true, nil
}
