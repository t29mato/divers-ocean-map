package service

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"

	"hello-world/model"
)

// DynamoDBServiceImpl ...
type DynamoDBServiceImpl struct {
	endpoint  string
	tableName string
	dynamoDB  *dynamodb.DynamoDB
}

// DynamoDBService ...
type DynamoDBService interface {
	NewDynamoDBService() *DynamoDBServiceImpl
	Create(*model.Ocean) error
}

// NewDynamoDBService ...
func NewDynamoDBService() *DynamoDBServiceImpl {
	fmt.Println("DYNAMODB_ENDPOINT:", os.Getenv("DYNAMODB_ENDPOINT"))
	fmt.Println("DYNAMODB_TABLE_NAME:", os.Getenv("DYNAMODB_TABLE_NAME"))
	s := &DynamoDBServiceImpl{
		endpoint:  os.Getenv("DYNAMODB_ENDPOINT"),
		tableName: os.Getenv("DYNAMODB_TABLE_NAME"),
		dynamoDB:  nil,
	}

	// DynamoDBの設定
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-1")
	if len(s.endpoint) > 0 {
		config = config.WithEndpoint(s.endpoint)
	}

	s.dynamoDB = dynamodb.New(sess, config)
	fmt.Println("s.tableName:", s.tableName)
	return s
}

// Create ...
func (s *DynamoDBServiceImpl) Create(ocean *model.Ocean) error {
	fmt.Println("s2.tableName:", s.tableName)
	resp, err := s.dynamoDB.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(uuid.New().String()),
			},
		},
	})
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}
