## 图片比例缩放工具

### Installation
    go get -u github.com/260by/tools/image

### Quick start
	package main

	import (
		"fmt"
		"github.com/260by/tools/image"
	)

	func main()  {
		// 根据指定宽度按比例缩放图片
		err := image.Scale("./images/1.jpeg", "./images/1-100.jpeg", 100,)
		if err != nil {
			panic(err)
		}

		// 根据指定宽度和高度缩放图片
		err = image.Thumbnail("./images/1.jpeg", "./images/1-100.jpeg", 100, 200)
		if err != nil {
			panic(err)
		}

		// 通过x,y坐标剪切图片，原图为440x660
		err = image.Cut("./images/1.jpeg", "./images/1-440x660.jpeg", 0, 0, 440, 610)
		if err != nil {
			panic(err)
		}

		w, h, err := image.GetImgWidthHeight("http://a-cdn.vogued.cn/images/20181130/78673733b0504f0e86b1ccc8118299af.jpeg")
		if err != nil {
			panic(err)
		}

		fmt.Println(w, h)
	}
