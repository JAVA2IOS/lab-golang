package model

type XlsxFile struct {
	Name string `json:"name"`
	Size int64 `json:"size"`
	CreateTime int64 `json:"create"`
}