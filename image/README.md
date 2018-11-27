## 图片比例缩放工具

### Quick start
    go get -u github.com/260by/tools/image

### Code example
    package main
    
    import (
    	"flag"
    	"github.com/260by/tools/image"
    )
    
    func main()  {
    	srcFile := flag.String("s", "", "Source file path")
    	dstFile := flag.String("d", "", "Target file path")
    	newWidth := flag.Int("w", 0, "New image file with")
    	flag.Parse()
    
    	image.Scale(*srcFile, *dstFile, *newWidth)
    }
