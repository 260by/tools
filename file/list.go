package file

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

// List 获取指定目录含子目录下的所有文件, 返回列表
func List(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, List(filepath.Join(path, file.Name()))...)
		} else {
			fileList = append(fileList, filepath.Join(path, file.Name()))
		}
	}

	return fileList
}
