package sftp

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

func Put(user, password, key, host, port, src, dst string) bool {
	sftpClient, err := Connect(user, password, key, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	// var remoteFilePath = file
	// var localDir = dst

	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(src)
	dstFile, err := sftpClient.Create(path.Join(dst, remoteFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	f, err := ioutil.ReadAll(srcFile)
	if err != nil {
		log.Fatal(err)
	}

	dstFile.Write(f)
	// fmt.Printf("%s Upload file to remote finished!", src)


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
	return true
}
