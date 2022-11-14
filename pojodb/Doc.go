package pojodb

import "Golangapi/config"

type Doc struct {
	Doc_name       string `json:"name"`
	User_id        string `json:"user_id"`
	Doc_id         string `json:"doc_id"`
	Date_created   string `json:"date_created"`
	S3_bucket_path string `json:"s3_bucket_path"`
}

func GetAlldocs_db() []Doc {
	var docs []Doc
	config.DB.Find(&docs)
	return docs
}

func Getdoc_db(docID string) Doc {
	var doc Doc
	config.DB.Where("doc_id=?", docID).First(&doc)
	return doc
}

func Postdocs_db(doc Doc) Doc {
	config.DB.Create(&doc)
	return doc
}

func DeleteDoc(docID string) {
	doc := Doc{}
	config.DB.Where("doc_id=?", docID).Delete(&doc)
}
