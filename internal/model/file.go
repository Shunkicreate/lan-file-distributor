package model

type File struct {
    Name string `json:"name"`
    Path string `json:"path"`
    Size int64  `json:"size"`
}
