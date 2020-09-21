package service

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

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
	fmt.Println("s.dynamoDB:", s.dynamoDB)
	return s
}

// Create ...
func (s *DynamoDBServiceImpl) Create(ocean *model.Ocean) error {
	av, err := dynamodbattribute.MarshalMap(ocean)
	if err != nil {
		return err
	}

	fmt.Println(av)

	putItem := &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      av,
	}
	_, err = s.dynamoDB.PutItem(putItem)
	if err != nil {
		return err
	}
	return nil
}
