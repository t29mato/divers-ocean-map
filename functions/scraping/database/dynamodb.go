package database

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"scraping/model"
)

// DynamoDBDatabaseImpl ...
type DynamoDBDatabaseImpl struct {
	tableName string
	dynamoDB  *dynamodb.DynamoDB
}

// DynamoDBDatabase ...
type DynamoDBDatabase interface {
	CreateIfNotExist(*model.Ocean) error
}

// NewDynamoDBDatabase ...
func NewDynamoDBDatabase() *DynamoDBDatabaseImpl {
	if os.Getenv("ENV") == "local" {
		endpoint := os.Getenv("DYNAMODB_ENDPOINT")
		s := &DynamoDBDatabaseImpl{
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
	return &DynamoDBDatabaseImpl{
		tableName: os.Getenv("DYNAMODB_TABLE_NAME"),
		dynamoDB:  dynamodb.New(session.Must(session.NewSession())),
	}
}

// CreateIfNotExist パーティションキーとレンジキーの両方が存在しない場合のみ新規レコード作成
func (s *DynamoDBDatabaseImpl) CreateIfNotExist(ocean *model.Ocean) error {
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
