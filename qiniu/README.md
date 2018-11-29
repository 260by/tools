## 七牛云存储上传工具

### Install 
    go get -u github.com/260by/tools/qiniu

### Quick Start
    package main
    
    import (
    	"encoding/json"
    	"fmt"
    	"flag"
    	"log"
    	"github.com/260by/tools/qiniu"
    )
    
    func main()  {
    	source := flag.String("src", "", "Upload local file or dir path")
    	dst := flag.String("dst", "", "Qiniu bucket path")
    	flag.Parse()
    
    	qiniu := qiniu.Qiniu{
    		Bucket: "bucket-name",
    		AccessKey: "access-key",
    		SecretKey: "secret-key",
    	}
    
    	results, err := qiniu.Upload(*source, *dst)
    	if err != nil {
    		log.Fatalln(err)
    	}
    
    	for _, result := range results {
    		rJSON, err := json.Marshal(result)
    		if err != nil {
    			log.Fatalln(err)
    		}
    		fmt.Println(string(rJSON))
    	}
    }
