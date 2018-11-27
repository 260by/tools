package image

import (
	"code.google.com/p/graphics-go/graphics"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

// LoadImage 读取文件
func LoadImage(path string) (img image.Image, extName string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()
	// 解码图片， extName获取图片文件扩展名
	img, extName, err = image.Decode(file)
	return
}

// SaveImage 保存文件
func SaveImage(path, extName string, img image.Image) (err error) {
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	// 根据文件扩展名编码图片
	switch extName {
	case "jpeg":
		return jpeg.Encode(imgfile, img, &jpeg.Options{Quality: 86})
	case "png":
		return png.Encode(imgfile, img)
	case "gif":
		return gif.Encode(imgfile, img, &gif.Options{})
	default:
		return errors.New("ERROR FORMAT")
	}
}

// Scale 缩放图片
func Scale(srcFile, dstFile string, newWidth int) (err error) {
	srcImg, extName, err := LoadImage(srcFile)
	if err != nil {
		return
	}
	bounds := srcImg.Bounds()
	dx := bounds.Dx()	// 原图片宽度
	dy := bounds.Dy()	// 原图片高度
	var dstImg *image.RGBA
	if newWidth == 0 {
		dstImg = image.NewRGBA(image.Rect(0, 0, dx, dy))
	} else {
		dstImg = image.NewRGBA(image.Rect(0, 0, newWidth, newWidth*dy/dx))  // 指定宽度按比例缩放
	}

	err = graphics.Scale(dstImg, srcImg)
	if err != nil {
		return
	}
	return SaveImage(dstFile, extName, dstImg)
}
