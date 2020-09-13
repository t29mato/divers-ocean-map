package service

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
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
}

// NewDynamoDBService ...
func (s *DynamoDBServiceImpl) NewDynamoDBService() *DynamoDBServiceImpl {
	// 環境変数の指定
	s.endpoint = os.Getenv("DYNAMODB_ENDPOINT")
	s.tableName = os.Getenv("DYNAMODB_TABLE_NAME")

	// DynamoDBの設定
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-1")
	if len(s.endpoint) > 0 {
		config = config.WithEndpoint(s.endpoint)
	}

	s.dynamoDB = dynamodb.New(sess, config)
	return s
}

// Create ...
func (s *DynamoDBServiceImpl) Create() error {
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
