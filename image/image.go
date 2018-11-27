package image

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"image/gif"
	"os"
	"code.google.com/p/graphics-go/graphics"
)

// LoadImage 读取文件
func LoadImage(path string) (img image.Image, extName string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()
	img, extName, err = image.Decode(file)
	return
}

// SaveImage 保存文件
func SaveImage(path, extName string, img image.Image) (err error) {
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	switch extName {
	case "jpeg":
		return jpeg.Encode(imgfile, img, &jpeg.Options{Quality: 100})
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
	dx := bounds.Dx()
	dy := bounds.Dy()
	dstImg := image.NewRGBA(image.Rect(0, 0, newWidth, newWidth*dy/dx))
	err = graphics.Scale(dstImg, srcImg)
	if err != nil {
		return
	}
	return SaveImage(dstFile, extName, dstImg)
}