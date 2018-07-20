package sftp

import (
	"io/ioutil"
	"os"
	"path"
	"github.com/pkg/sftp"
)

func Put(sftpClient *sftp.Client , src, dst string) (result bool, err error) {
	// sftpClient, err := Connect(user, password, key, host, port)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer sftpClient.Close()

	// var remoteFilePath = file
	// var localDir = dst

	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(src)
	dstFile, err := sftpClient.Create(path.Join(dst, remoteFileName))
	if err != nil {
		return
	}
	defer dstFile.Close()

	f, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return
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
	return
}
