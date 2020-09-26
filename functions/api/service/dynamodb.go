package service

import (
	"api/logging"
	"api/model"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBServiceImpl ...
type DynamoDBServiceImpl struct {
	tableName string
	dynamoDB  *dynamodb.DynamoDB
	logging   *logging.OceanLoggingImpl
}

// DynamoDBService ...
type DynamoDBService interface {
	CreateIfNotExist(*model.Ocean) error
}

// NewDynamoDBService ...
func NewDynamoDBService(logging *logging.OceanLoggingImpl) *DynamoDBServiceImpl {
	if os.Getenv("ENV") == "local" {
		endpoint := os.Getenv("DYNAMODB_ENDPOINT")
		s := &DynamoDBServiceImpl{
			tableName: os.Getenv("DYNAMODB_TABLE_NAME"),
			dynamoDB:  nil,
			logging:   logging,
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

// FetchLatestOcean 指定されたダイビングポイントの最新の海況情報を返す
func (s *DynamoDBServiceImpl) FetchLatestOcean(locationName string) (*model.Ocean, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(s.tableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":LocationName": {S: aws.String(locationName)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#LocationName": aws.String("LocationName"),
		},
		KeyConditionExpression: aws.String("#LocationName = :LocationName"),
		Limit:                  aws.Int64(1),
		ScanIndexForward:       aws.Bool(false), // 最新の日付順にするため (デフォルトだと古い順)
	}
	output, err := s.dynamoDB.Query(queryInput)
	fmt.Println(output)

	var ocean model.Ocean

	// TODO: 0件の場合にフロントで処理できるように、なんらか手を考える
	if len(output.Items) == 0 {
		s.logging.Info("該当するレコードは存在しません。, locationName=", locationName)
		return nil, nil
	}
	err = dynamodbattribute.UnmarshalMap(output.Items[0], &ocean)
	if err != nil {
		return nil, err
	}
	return &ocean, nil
}

// FetchAllLatestOceans 全てのダイビングポイントの最新の海況情報を返す
func (s *DynamoDBServiceImpl) FetchAllLatestOceans() ([]*model.Ocean, error) {
	oceanNameList := []string{
		"izu-ocean-park",
		"ukishima-in-tiba-katsuyama",
	}

	var oceans []*model.Ocean
	for _, name := range oceanNameList {
		ocean, err := s.FetchLatestOcean(name)
		if err != nil {
			s.logging.Info("取得に失敗, locationName:", name)
			continue
		}
		oceans = append(oceans, ocean)
	}

	return oceans, nil
}
