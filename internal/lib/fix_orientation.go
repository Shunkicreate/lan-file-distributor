package lib

import (
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"io"
)

// ExtractOrientation reads the EXIF orientation tag from the image file.
// If EXIF or orientation info is missing, it returns 1 (normal orientation).
func ExtractOrientation(r io.Reader) (int, error) {
	x, err := exif.Decode(r)
	if err != nil {
		// No EXIF: assume normal orientation
		return 1, nil
	}
	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return 1, nil
	}
	orientation, err := tag.Int(0)
	if err != nil {
		return 1, nil
	}
	return orientation, nil
}

// RotateByOrientation rotates the image according to the EXIF orientation tag value.
func RotateByOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 3:
		return imaging.Rotate180(img)
	case 6:
		return imaging.Rotate270(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}
