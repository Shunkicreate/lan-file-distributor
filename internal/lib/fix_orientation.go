// lib/fix_orientation.go
package lib

import (
	"github.com/rwcarlsen/goexif/exif"
	"github.com/disintegration/imaging"
	"image"
	"os"
)

func FixOrientation(file *os.File, img image.Image) (image.Image, error) {
	// Reset file pointer to beginning
	file.Seek(0, 0)

	exifData, err := exif.Decode(file)
	if err != nil {
		// Exifがない場合はそのまま返す
		return img, nil
	}

	orientationTag, err := exifData.Get(exif.Orientation)
	if err != nil {
		return img, nil
	}

	orientation, err := orientationTag.Int(0)
	if err != nil {
		return img, nil
	}

	switch orientation {
	case 3:
		return imaging.Rotate180(img), nil
	case 6:
		return imaging.Rotate270(img), nil
	case 8:
		return imaging.Rotate90(img), nil
	default:
		return img, nil
	}
}
