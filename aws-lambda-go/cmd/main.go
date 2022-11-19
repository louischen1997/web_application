package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) error {

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		Sendemails(snsRecord.Message)
	}
	return nil
}

func Sendemails(mess string) {
	res := strings.Split(mess, "/")
	res1 := res[0]
	res2 := res[1]
	res3 := res[2]
	res4 := res[3]
	emailaddress := "http://" + res4 + ":3000/v1/verifyUserEmail?email=" + res1 + "&token=" + res2
	from := mail.NewEmail("Web Application CSYE", "csye6225@em6674.prod.louisdomain6225.me")
	subject := "Verify User Account"
	to := mail.NewEmail("New User", res1)
	plainTextContent := "Verify User Account"
	htmlContent := "<strong>Please click the following link to verify your account<br /><br /><div>" + emailaddress + "</div></strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(res3)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
