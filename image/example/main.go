package main

import (
	"fmt"
	"github.com/260by/tools/image"
)

func main()  {
	// 根据指定宽度按比例缩放图片
	err := image.Scale("./images/1.jpeg", "./images/1-100.jpeg", 200,)
	if err != nil {
		panic(err)
	}

	// 根据指定宽度和高度缩放图片
	err = image.Thumbnail("./images/1.jpeg", "./images/1-236x354.jpeg", 236, 354)
	if err != nil {
		panic(err)
	}

	// 通过x,y坐标剪切图片，原图为440x660
	err = image.Cut("./images/1.jpeg", "./images/1-440x660.jpeg", 0, 0, 440, 640)
	if err != nil {
		panic(err)
	}

	// 获取图片宽度和高度
	w, h, err := image.GetImgWidthHeight("./images/1.jpeg")
	if err != nil {
		panic(err)
	}

	fmt.Println(w, h)
}