package contourmap

import (
	"image"
	"image/draw"
)

func imageToGray16(im image.Image) *image.Gray16 {
	dst := image.NewGray16(im.Bounds())
	draw.Draw(dst, im.Bounds(), im, image.ZP, draw.Src)
	return dst
}
