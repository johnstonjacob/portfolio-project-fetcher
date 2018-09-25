package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	region      string
	credentials string
	profile     string
)

func init() {
	region = os.Getenv("AWS_REGION")
}

func main() {
	lambda.Start(projectHandler)
}

func projectHandler() (events.APIGatewayProxyResponse, error) {
	b := body{}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("portfolio"),
	})

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &b.Projects)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	bjson, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}
	r := events.APIGatewayProxyResponse{Headers: map[string]string{}, IsBase64Encoded: false, StatusCode: 200, Body: string(bjson)}

	return r, nil
}

type response struct {
	Headers         map[string]string `json:"headers"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Body            body              `json:"body"`
}

type body struct {
	Projects []project `json:"projects"`
}

type project struct {
	Project      string `json:"project"`
	Brief        string `json:"brief"`
	Technologies string `json:"technologies"`
	GithubURL    string `json:"githubURL"`
}
