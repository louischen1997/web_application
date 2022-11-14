package mdldb

type Doc struct {
	Doc_name       string `gorm:"type:varchar(45)"`
	Doc_ID         string `gorm:"type:varchar(45);primary_key"`
	User_ID        string `gorm:"type:varchar(45)"`
	Date_created   string `gorm:"type:varchar(100)"`
	S3_bucket_path string `gorm:"type:varchar(150)"`
}

type Docs []Doc
