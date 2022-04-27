package tools

import (
	"China_Telecom_Monitor/configs"
	"bytes"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

// 图片加水印
func Watermark(imgByte []byte, text, filePath string) {

	reader := bytes.NewReader(imgByte)
	pngimg, err := png.Decode(reader)
	if err != nil {
		configs.Logger.Error(err)
		return
	}

	img := image.NewNRGBA(pngimg.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, pngimg.At(x, y))
		}
	}

	font, err := freetype.ParseFont(ResourceSimHeiTtf.StaticContent)
	if err != nil {
		configs.Logger.Error(err)
		return
	}

	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(75)
	f.SetClip(pngimg.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(color.RGBA{R: 119, G: 136, B: 153, A: 255}))

	pt := freetype.Pt(img.Bounds().Dx()-1000, img.Bounds().Dy()-25)
	_, err = f.DrawString(text, pt)

	//draw.Draw(img,jpgimg.Bounds(),jpgimg,image.ZP,draw.Over)

	//保存到新文件中
	newfile, err := Create(filePath)
	if err != nil {
		configs.Logger.Error(err)
		return
	}
	defer newfile.Close()

	err = jpeg.Encode(newfile, img, &jpeg.Options{100})
	if err != nil {
		configs.Logger.Error(err)
	}
}
