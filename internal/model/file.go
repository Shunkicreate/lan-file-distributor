package model

import "image"

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type ImageFile struct {
	Image    image.Image `json:"-"`      // JSONシリアライズから除外
	Path     string     `json:"path"`
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Data     []byte     `json:"data"`    // エンコードされた画像データ
}
