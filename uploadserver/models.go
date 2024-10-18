package uploadserver

type S3Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}
type Report struct {
	FileUrl    string `bson:"fileUrl"`
	FileName   string `bson:"fileName"`
	UploadedAt string `bson:"uploadedAt"`
}

type User struct {
	UserId    string   `bson:"userId"`
	Name      string   `bson:"name"`
	Email     string   `bson:"email"`
	Gender    string   `bson:"gender"`
	Age       int      `bson:"age"`
	Reports   []Report `bson:"reports"`
	CreatedAt string   `bson:"createdAt"`
}
