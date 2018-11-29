package qiniu

import (
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Qiniu 七牛账号信息
type Qiniu struct {
	Bucket    string
	AccessKey string
	SecretKey string
}

// Result 返回信息结构体
type Result struct {
	Key        string
	Hash       string
	Fsize      int
	URL        string
	SourceFile string
	Status     string
}

// GetBucketDomain 根据bucket获取域名
func (q *Qiniu) GetBucketDomain() (domain string, err error) {
	client := &http.Client{}
	url := fmt.Sprintf("http://api.qiniu.com/v6/domain/list?tbl=%s", q.Bucket)
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	token, err := mac.SignRequest(reqest)
	if err != nil {
		return
	}

	// 构造header
	auth := fmt.Sprintf("QBox %s", token)
	reqest.Header.Add("Authorization", auth)
	reqest.Header.Add("Content-Type", "application/json")

	response, err := client.Do(reqest)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	d := strings.Trim(string(body), "[]")
	l := strings.Split(d, ",")
	for _, k := range l {
		// 判断域名是否七牛测试域名
		result, err := regexp.Match("clouddn.com", []byte(k))
		if err != nil {
			continue
		}
		if !result {
			domain = k
		}
	}

	return strings.Trim(domain, "\""), nil
}

// Upload 上传
func (q *Qiniu) Upload(src, dst string) (results []Result, err error) {
	putPolicy := storage.PutPolicy{
		Scope:      q.Bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize)}`,
	}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := Result{}

	// 如果src是目录则获取目录里所有文件
	var fileList []string
	f, err := os.Stat(src)
	if err != nil {
		return
	}
	if f.IsDir() {
		fileList = getFileList(src)
	} else {
		fileList = append(fileList, src)
	}

	uploadDir := strings.Trim(dst, "/")
	uploadDir += "/"

	baseDomain, err := q.GetBucketDomain()
	if err != nil {
		return
	}

	for _, file := range fileList {
		var key string
		if f.IsDir() {
			relativePath := strings.TrimPrefix(strings.TrimPrefix(file, src), "/")
			key = uploadDir + relativePath
		} else {
			key = uploadDir + filepath.Base(file)
		}

		err := formUploader.PutFile(context.Background(), &ret, upToken, key, file, nil)
		if err != nil {
			results = append(results, Result{SourceFile: file, Status: "Failed"})
		} else {
			results = append(results, Result{
				Key:        ret.Key,
				Hash:       ret.Hash,
				Fsize:      ret.Fsize,
				URL:        baseDomain + "/" + ret.Key,
				SourceFile: file,
				Status:     "Success",
			})
		}

	}

	return results, nil
}

// List 获取指定目录含子目录下的所有文件, 返回列表
func getFileList(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, getFileList(filepath.Join(path, file.Name()))...)
		} else {
			fileList = append(fileList, filepath.Join(path, file.Name()))
		}
	}

	return fileList
}
