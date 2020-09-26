package service

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"scraping/model"
)

// DynamoDBServiceImpl ...
type DynamoDBServiceImpl struct {
	tableName string
	dynamoDB  *dynamodb.DynamoDB
}

// DynamoDBService ...
type DynamoDBService interface {
	CreateIfNotExist(*model.Ocean) error
}

// NewDynamoDBService ...
func NewDynamoDBService() *DynamoDBServiceImpl {
	if os.Getenv("ENV") == "local" {
		endpoint := os.Getenv("DYNAMODB_ENDPOINT")
		s := &DynamoDBServiceImpl{
			tableName: os.Getenv("DYNAMODB_TABLE_NAME"),
			dynamoDB:  nil,
		}
		// DynamoDBの設定
		sess := session.Must(session.NewSession())
		config := aws.NewConfig().WithRegion("ap-northeast-1")
		if len(endpoint) > 0 {
			config = config.WithEndpoint(endpoint)
		}

		s.dynamoDB = dynamodb.New(sess, config)
		return s
	}

	// AWS上では、endpointなしで、自動で解決してくれるため、endpointの設定なし
	return &DynamoDBServiceImpl{
		tableName: os.Getenv("DYNAMODB_TABLE_NAME"),
		dynamoDB:  dynamodb.New(session.Must(session.NewSession())),
	}
}

// CreateIfNotExist パーティションキーとレンジキーの両方が存在しない場合のみ新規レコード作成
func (s *DynamoDBServiceImpl) CreateIfNotExist(ocean *model.Ocean) error {
	av, err := dynamodbattribute.MarshalMap(ocean)
	if err != nil {
		return err
	}

	putItem := &dynamodb.PutItemInput{
		TableName:           aws.String(s.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(LocationName) AND attribute_not_exists(MeasuredTime)"),
	}
	_, err = s.dynamoDB.PutItem(putItem)
	if err != nil {
		return err
	}
	return nil
}
