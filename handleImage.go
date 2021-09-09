// Package services /**
package services

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/golang/freetype"
	"github.com/nfnt/resize"
)

type ParamsImage struct {
	Path       string
	Width      int
	Haight     int
	ImageModel image.Image
}

type HandleImage struct {
	ParamsImage
	Images []EleImage
	Tests  []EleText
}

type EleImage struct {
	ParamsImage
	X int
	Y int
}

type EleText struct {
	Content    string
	X          int
	Y          int
	Size       int
	Color      string
	ImageModel image.Image
}

func getImageDecode(imgPath string) (image.Image, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	imgDecode, _, decodeErr := image.Decode(file)
	if decodeErr != nil {
		return nil, decodeErr
	}
	defer file.Close()
	return imgDecode, nil
}

func (eleText *EleText) HandleTest() error {
	img := image.NewNRGBA(image.Rect(0, 0, 200, 200))
	fontBytes, err := ioutil.ReadFile("可爱萌萌字体ttf.ttf")
	if err != nil {
		log.Println(err)
	}
	//载入字体数据
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}
	f := freetype.NewContext()
	//设置分辨率
	f.SetDPI(72)
	//设置字体
	f.SetFont(font)
	//设置尺寸
	f.SetFontSize(float64(eleText.Size))
	f.SetClip(img.Bounds())
	//设置输出的图片
	f.SetDst(img)
	//设置字体颜色(红色)
	f.SetSrc(image.NewUniform(color.RGBA{255, 0, 0, 255}))
	//设置字体的位置
	pt := freetype.Pt(40, 40+int(f.PointToFixed(26))>>8)

	_, err = f.DrawString(eleText.Content, pt)
	if err != nil {
		return err
	}
	eleText.ImageModel = img
	return nil
}

func (pImage *ParamsImage) CuttingImage() error {
	if pImage.ImageModel != nil {
		return nil
	}
	imgDecode, err := getImageDecode(pImage.Path)
	if err != nil {
		return err
	}
	if pImage.Width == 0 && pImage.Haight == 0 {
		pImage.ImageModel = imgDecode
		return nil
	}
	imageBounds := imgDecode.Bounds()
	if pImage.Width == imageBounds.Max.X && pImage.Haight == imageBounds.Max.Y {
		pImage.ImageModel = imgDecode
		return nil
	}
	afterImage := resize.Resize(uint(pImage.Width), uint(pImage.Haight), imgDecode, resize.Lanczos2)
	pImage.ImageModel = afterImage
	return nil
}

// func Crea
func (bgImage *HandleImage) ComposeImage() error {
	if err := bgImage.CuttingImage(); err != nil {
		return err
	}
	resImage := image.NewRGBA(bgImage.ImageModel.Bounds())
	draw.Draw(resImage, resImage.Bounds(), bgImage.ImageModel, image.Point{X: 0, Y: 0}, draw.Src)
	bgImage.composeImageForImages(resImage)
	bgImage.composeImagesForTexts(resImage)
	resFile, _ := os.Create(path.Base(fmt.Sprintf("TEST%d.jpeg", time.Now().UnixNano())))
	_ = jpeg.Encode(resFile, resImage, &jpeg.Options{Quality: 100})
	defer resFile.Close()
	return nil
}

func (tImage *HandleImage) composeImageForImages(resImage draw.Image) error {
	if len(tImage.Images) == 0 {
		return nil
	}
	for _, v := range tImage.Images {
		err := v.CuttingImage()
		if err != nil {
			return nil
		}
		offset := image.Pt(v.X, v.Y)
		draw.Draw(resImage, v.ImageModel.Bounds().Add(offset), v.ImageModel, image.Point{X: 0, Y: 0}, draw.Over)
	}
	return nil
}

func (text *HandleImage) composeImagesForTexts(resImage draw.Image) error {
	if len(text.Images) == 0 {
		return nil
	}
	for _, v := range text.Tests {
		err := v.HandleTest()
		if err != nil {
			return nil
		}
		offset := image.Pt(v.X, v.Y)
		draw.Draw(resImage, v.ImageModel.Bounds().Add(offset), v.ImageModel, image.Point{X: 0, Y: 0}, draw.Over)
	}
	return nil
}
