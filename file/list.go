package file

import (
	"io/ioutil"
	"path/filepath"
	"log"
)

// 获取指定目录包含子目录下的所有文件
// 返回文件列表
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