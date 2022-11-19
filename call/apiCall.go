package call

import (
	"Golangapi/pojo"
	"fmt"
	"log"
	"net/http"
	"os"

	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Item struct {
	EMAIL_KEY string                   `json:"EMAIL_KEY"`
	TOKEN     string                   `json:"TOKEN"`
	TTLATT    dynamodbattribute.Number `json:"TTLATT"`
}

var hasAuth bool = false
var tmpuserID string = ""

var Logger *log.Logger



func MyLogger() {
	os.Create("csye6225.log")
	var logPath = "csye6225.log"
	var errFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error", err)
	}
	Logger = log.New(errFile, "log: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func HandleMetricCounter(clientName string) {
	config := &statsd.ClientConfig{
		Address: "127.0.0.1:8125",
		Prefix:  "statsd-client ",
	}
	client, err := statsd.NewClientWithConfig(config)
	if err != nil {
		Logger.Print(err)
	}
	client.Inc(clientName, 1, 1.0)
	defer client.Close()
}

func GetAllUsers(c *gin.Context) {
	users := pojo.GetAllUsers_db()
	c.JSON(200, users)
}

func GetAllUsersToCheck(email string) bool {
	users := pojo.GetAllUsers_db()
	for index, element := range users {
		if element.Username == email {
			return false
		}
		index += 1
	}
	return true
}

func GetUsers(c *gin.Context) {
	ep := os.Getenv("endpoint")
	user := pojo.GetUsers_db(c.Param("id"))
	HandleMetricCounter("Get_Users")

	if !hasAuth {
		c.Status(401)
		Logger.Print("Get Request, Authenticate Fail, Endpoint: " + ep + ":3000/v1/account/" + user.ID)
		return
	}
	if user.VerifyType == "0" {
		c.Status(401)
		Logger.Print("Get Request, Authenticate Fail->unverified, Endpoint: " + ep + ":3000/v1/account/" + user.ID)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"first_name":      user.First_name,
		"last_name":       user.Last_name,
		"username":        user.Username,
		"id":              user.ID,
		"account_created": user.Account_created,
		"account_updated": user.Account_updated,
		"verify_status":   "verified",
	})
	Logger.Print("Get Request, Get User ID: " + user.ID + ", Endpoint: " + ep + ":3000/v1/account/" + user.ID)
	hasAuth = false

}

func PostUsers(c *gin.Context) {
	//sending()
	HandleMetricCounter("Post_Users")
	ep := os.Getenv("endpoint")
	user := pojo.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.Status(400)
		Logger.Print("Post Request, Bad request 400, Endpoint: " + ep + ":3000/v1/account")
		return
	}

	if !GetAllUsersToCheck(user.Username) {
		c.Status(400)
		Logger.Print("Post Request, Bad Request 400, Endpoint: " + ep + ":3000/v1/account")

		return
	}
	if user.ID != "" || user.Account_created != "" || user.Account_updated != "" {
		c.Status(400)
		Logger.Print("Post Request, Bad Request 400, Endpoint: " + ep + ":3000/v1/account")

		return
	}

	hashPass := []byte(user.Password)
	user.Password = HashAndSalt(hashPass)

	id := uuid.New().String()
	now := time.Now().String()

	user.ID = id
	user.Account_created = now
	user.Account_updated = now
	user.VerifyType = "0"

	pojo.PostUsers_db(user)
	c.JSON(201, gin.H{
		"first_name":      user.First_name,
		"last_name":       user.Last_name,
		"username":        user.Username,
		"id":              user.ID,
		"account_created": user.Account_created,
		"account_updated": user.Account_updated,
		"verify_status":   "unverified",
	})
	Logger.Print("Post Request, Post New User: " + user.First_name + " " + user.Last_name + ", Endpoint: " + ep + ":3000/v1/account")

	oneTimeToken := uuid.New().String()

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)
	sendgridkey := os.Getenv("sendgrid")

	mess := user.Username + "/" + oneTimeToken + "/" + sendgridkey + "/" + ep

	topicarn := os.Getenv("snstopic")
	result, errr := svc.Publish(&sns.PublishInput{
		Message:  &mess,
		TopicArn: &topicarn,
	})

	if errr != nil {
		fmt.Println(errr.Error())
		os.Exit(1)
	}

	fmt.Println(*result.MessageId)

	//sendemails()

	item := Item{
		EMAIL_KEY: user.Username,
		TOKEN:     oneTimeToken,
		TTLATT:    dynamodbattribute.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	svc2 := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		Logger.Print(err)
	}

	tableName := os.Getenv("dynamotablename")

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc2.PutItem(input)
	if err != nil {
		Logger.Print(err)
	}

}

func DeleteUser(c *gin.Context) {
	pojo.DeleteUser(c.Param("id"))
	c.Status(http.StatusOK)
}

func DeleteAllUsers(c *gin.Context) {
	users := pojo.GetAllUsers_db()
	for index, element := range users {
		pojo.DeleteUser(element.ID)
		index += 1
	}
	c.JSON(200, "DELETE ALL")
}

func PutUser(c *gin.Context) {
	HandleMetricCounter("Put_Users")
	ep := os.Getenv("endpoint")
	user := pojo.User{}
	tmp_user := pojo.GetUsers_db(c.Param("id"))
	if !hasAuth {
		Logger.Print("Put Request, Authenticate Fail, Endpoint: " + ep + ":3000/v1/account/" + tmp_user.ID)
		return
	}
	if tmp_user.VerifyType == "0" {
		c.Status(401)
		Logger.Print("Put Request, Authenticate Fail->unverified, Endpoint: " + ep + ":3000/v1/account/" + user.ID)
		return
	}

	now := time.Now().String()

	err := c.BindJSON(&user)
	if err != nil {
		Logger.Print("Put Request, Bad Request 400, Endpoint: " + ep + ":3000/v1/account/" + tmp_user.ID)
		c.Status(400)
		return
	}
	if user.ID != "" || user.Account_created != "" || user.Account_updated != "" || user.Username != tmp_user.Username {
		c.Status(400)
		Logger.Print("Put Request, Bad Request 400, Endpoint: " + ep + ":3000/v1/account/" + tmp_user.ID)
		return
	}
	hashPass := []byte(user.Password)
	user.Password = HashAndSalt(hashPass)
	user.Account_updated = now
	pojo.UpdateUser(c.Param("id"), user)
	hasAuth = false
	c.Status(204)
	Logger.Print("Put Request, Update User ID: " + user.ID + " Status Code 204 Update Successful, Endpoint: " + ep + ":3000/v1/account/" + tmp_user.ID)

}

func HandleBA(c *gin.Context) {
	var username_ba, password_ba string
	users := pojo.GetAllUsers_db()
	username_ba, password_ba, hasAuth = c.Request.BasicAuth()
	if len(users) == 0 {
		Logger.Print("No User In The List, Status Code 401 Unauthorized")
		c.Status(401)
		hasAuth = false
		return
	}
	for index, element := range users {
		if index == len(users)-1 && element.ID != c.Param("id") {
			Logger.Print("Wrong User ID, Status Code 401 Unauthorized")
			c.Status(401)
			hasAuth = false
			return
		}
		if element.ID == c.Param("id") {
			break
		}
	}
	tmpHash := pojo.GetUsers_db_Pass(c.Param("id"))
	if username_ba == pojo.GetUsers_db_Username(c.Param("id")) && CheckPasswordHash(password_ba, tmpHash) {

		Logger.Print("Basic Auth Successful")
		return
	}
	hasAuth = false
	c.Status(403)
	Logger.Print("Authenticate Fail, Status Code 403 Forbidden")
}

func HandleBA_doc(c *gin.Context) {
	var username_ba, password_ba string
	var tmpid = ""
	users := pojo.GetAllUsers_db()
	username_ba, password_ba, hasAuth = c.Request.BasicAuth()
	if len(users) == 0 {
		Logger.Print("No User In The List, Status Code 401 Unauthorized")
		c.Status(401)
		hasAuth = false
		return
	}
	for index, element := range users {
		if index == len(users)-1 && element.Username != username_ba {
			c.Status(401)
			hasAuth = false
			Logger.Print("Wrong User ID, Status Code 401 Unauthorized")
			return
		}
		if element.Username == username_ba {
			tmpid = element.ID
			break
		}
	}
	tmpHash := pojo.GetUsers_db_Pass(tmpid)
	if username_ba == pojo.GetUsers_db_Username(tmpid) && CheckPasswordHash(password_ba, tmpHash) {
		Logger.Print("Basic Auth Successful")

		user := pojo.GetUsers_db(tmpid)
		if user.VerifyType == "0" {
			c.Status(401)
			Logger.Print("Authenticate Fail->unverified,401")
			hasAuth = false
			return
		}

		tmpuserID = tmpid
		return
	}
	hasAuth = false
	c.Status(401)
	Logger.Print("Status code 401 Unauthorized")

}

func HashAndSalt(pass []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func VerifyEmail(c *gin.Context) {
	tableName := os.Getenv("dynamotablename")
	email := c.Request.URL.Query().Get("email")
	token := c.Request.URL.Query().Get("token")
	Logger.Print("get from url")
	Logger.Print(email)
	Logger.Print(token)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	filt := expression.Name("EMAIL_KEY").Equal(expression.Value(email))
	proj := expression.NamesList(expression.Name("TTLATT"), expression.Name("EMAIL_KEY"), expression.Name("TOKEN"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		Logger.Print("Expression Error")
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		Logger.Print("Params Error")
	}

	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			Logger.Print(err)
		}

		Logger.Print("FIND ITEM")
		Logger.Print(item.EMAIL_KEY)
		Logger.Print(item.TOKEN)
		Logger.Print(item.TTLATT)

		tmptime := dynamodbattribute.Number(strconv.FormatInt(time.Now().Unix(), 10)).String()
		intVar1, err1 := strconv.Atoi(tmptime)
		if err1 != nil {
			Logger.Print("Convert int error")
		}
		intVar2, err2 := strconv.Atoi(string(item.TTLATT))
		if err2 != nil {
			Logger.Print("Convert int error")
		}
		Logger.Print("This is TTL DEBUG")
		Logger.Print(intVar1)
		Logger.Print(intVar2)
		Logger.Print(intVar1 - intVar2)

		if intVar1-intVar2 > 300 {
			Logger.Print("Verify Link Expired")
			c.Status(400)
			return
		}

		if token != item.TOKEN {
			c.Status(400)
			return
		} else {
			users := pojo.GetAllUsers_db()
			for index, element := range users {
				if element.Username == item.EMAIL_KEY {
					element.VerifyType = "1"
					pojo.UpdateUser(element.ID, element)
					return
				}
				index += 1
			}

		}

	}
	c.Status(400)
}
