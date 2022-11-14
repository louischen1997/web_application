package call

import (
	"Golangapi/pojodb"
	"context"
	"os"
	"time"

	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAlldocs(c *gin.Context) {
	HandleMetricCounter("Get_All_docs")
	ep := os.Getenv("endpoint")

	if !hasAuth {
		Logger.Print("Get Request, Authenticate Fail, Endpoint: " + ep + ":3000/v1/documents")
		return
	}
	c.JSON(200, Getall_id_docs(tmpuserID))

	Logger.Print("Get Request, Get All Docs From User ID: " + tmpuserID + ", 200 OK, Endpoint: " + ep + ":3000/v1/documents")
	tmpuserID = ""
	hasAuth = false

}

func Getdoc(c *gin.Context) {
	HandleMetricCounter("Get_doc")
	ep := os.Getenv("endpoint")

	if !hasAuth {
		Logger.Print("Get Request, Authenticate Fail, Endpoint: " + ep + ":3000/v1/documents/" + c.Param("id"))
		return
	}
	doc := pojodb.Getdoc_db(c.Param("id"))
	if doc.User_id == tmpuserID {
		c.JSON(http.StatusOK, doc)
		Logger.Print("Get Request, Get Doc ID" + tmpuserID + ", 200 OK, Endpoint: " + ep + ":3000/v1/documents/" + c.Param("id"))
	} else {
		c.Status(403)
		Logger.Print("Get Request, Wrong ID: 403 Forbidden, Endpoint: " + ep + ":3000/v1/documents/" + c.Param("id"))
	}
	tmpuserID = ""
	hasAuth = false

}

func Postdocs(c *gin.Context) {
	HandleMetricCounter("Post_doc")
	ep := os.Getenv("endpoint")
	if !hasAuth {
		Logger.Print("Post Request, Authenticate Fail, Endpoint: " + ep + ":3000/v1/documents")
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		Logger.Printf("Error: %v", err)
		return
	}
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	doc := pojodb.Doc{}

	file, err := c.FormFile("file")

	if err != nil {
		Logger.Print("Get file process error")
		return
	}

	id := uuid.New().String()
	now := time.Now().String()

	bkn := os.Getenv("bkn")

	doc.Doc_id = id
	doc.Date_created = now
	doc.User_id = tmpuserID
	doc.S3_bucket_path = "s3://" + bkn + "/" + id + "/" + file.Filename
	doc.Doc_name = file.Filename
	pojodb.Postdocs_db(doc)

	f, _ := file.Open()
	uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bkn),
		Key:    aws.String(id + "/" + file.Filename),
		Body:   f,
	})

	c.JSON(201, doc)
	Logger.Print("Post Request, Post New File To S3 Bucket, Filename: " + file.Filename + ", S3 Path: " + doc.S3_bucket_path + ", Endpoint: " + ep + ":3000/v1/documents")
	tmpuserID = ""
	hasAuth = false
}

func DeleteDoc(c *gin.Context) {
	HandleMetricCounter("Delete_doc")
	ep := os.Getenv("endpoint")

	if !hasAuth {
		Logger.Print("Delete Request, Authenticate Fail, Endpoint: " + ep + ":3000/v1/documents/" + c.Param("id"))
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		Logger.Printf("Error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)

	doc := pojodb.Getdoc_db(c.Param("id"))

	if doc.User_id == tmpuserID {
		pojodb.DeleteDoc(c.Param("id"))

		bkn := os.Getenv("bkn")

		_, newerr := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(bkn),
			Key:    aws.String(doc.Doc_id + "/" + doc.Doc_name),
		})
		if newerr != nil {
			panic("Couldn't Delete Object")
		}

		c.Status(204)
		Logger.Print("Delete Request, Successful Delete Object From " + doc.S3_bucket_path + ", Status code 204, Endpoint: " + ep + ":3000/v1/documents/" + c.Param("id"))
	} else {
		c.Status(404)
		Logger.Print("Delete Request, Record Not Found, 404 Not Found, Endpoint: " + ep + ":3000/v1/documents/" + c.Param("id"))

	}

	tmpuserID = ""
	hasAuth = false
}

func Getall_id_docs(id string) []pojodb.Doc {
	var docs []pojodb.Doc
	all_docs := pojodb.GetAlldocs_db()

	for index, element := range all_docs {
		if element.User_id == id {
			docs = append(docs, element)
		}
		index += 1
	}

	return docs
}

func DeleteAllDoc(c *gin.Context) {
	all_docs := pojodb.GetAlldocs_db()
	for index, element := range all_docs {
		if element.User_id == tmpuserID {
			pojodb.DeleteDoc(element.Doc_id)
		}
		index += 1
	}
	c.Status(200)
}
