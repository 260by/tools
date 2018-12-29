package image

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/260by/tools/image/graphics"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Thumbnail 按宽度和高度进行比例缩放
func Thumbnail(filePath string, savePath string, width, height int) (err error) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	src, filetype, err := LoadImage(filePath)
	if err != nil {
		return
	}
	err = graphics.Thumbnail(dst, src)
	if err != nil {
		return
	}
	err = SaveImage(savePath, dst, filetype)
	return
}

// Scale 按宽度进行比例缩放
func Scale(filePath string, savePath string, newWidth int) (err error) {
	srcImg, filetype, err := LoadImage(filePath)
	if err != nil {
		return
	}
	bound := srcImg.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()
	dstImg := image.NewRGBA(image.Rect(0, 0, newWidth, newWidth*dy/dx))
	// 产生缩略图,等比例缩放
	err = graphics.Scale(dstImg, srcImg)
	if err != nil {
		return
	}
	err = SaveImage(savePath, dstImg, filetype)
	if err != nil {
		return
	}
	return
}

// Cut 根据指定的x,y轴剪切图片,图片x,y坐标零点为左上角
func Cut(filePath string, savePath string, x0, y0, x1, y1 int) (err error) {
	src, filetype, err := LoadImage(filePath)
	if err != nil {
		return
	}

	out, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer out.Close()

	switch filetype {
	case "jpeg":
		img := src.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{Quality: 86})
	case "png":
		switch src.(type) {
		case *image.NRGBA:
			img := src.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := src.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := src.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	default:
		return errors.New("ERROR FORMAT")
	}

	return nil
}

// GetImgWidthHeight 获取图片的宽度和高度
func GetImgWidthHeight(filename string) (w, h int, err error) {
	img, _, err := LoadImage(filename)
	if err != nil {
		return
	}

	return img.Bounds().Dx(), img.Bounds().Dy(), nil
}

// LoadImage 根据文件名打开图片,并编码,返回编码对象和文件类型
func LoadImage(path string) (img image.Image, fileType string, err error) {
	u, _ := url.Parse(path)
	if u.Host != "" {
		response, err := http.Get(path)
		if err != nil {
			return nil, "", err
		}
		if response.StatusCode != 200 {
			return nil, "", fmt.Errorf("Error: %v", response.StatusCode)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, "", err
		}
		img, fileType, err = image.Decode(bytes.NewReader(body))
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, "", err
		}
		defer file.Close()
		img, fileType, err = image.Decode(file)
	}
	return
}

// SaveImage 将编码对象存入文件中
func SaveImage(path string, img *image.RGBA, fileType string) (err error) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return
	}
	// 根据文件扩展名编码图片
	switch fileType {
	case "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 86})
	case "png":
		return png.Encode(file, img)
	case "gif":
		return gif.Encode(file, img, &gif.Options{})
	default:
		return errors.New("ERROR FORMAT")
	}
}
