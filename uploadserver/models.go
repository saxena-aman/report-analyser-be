package uploadserver

type S3Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}
