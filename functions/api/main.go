package main

import (
	"api/logging"
	"api/service"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler ...
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logging := logging.NewOceanLoggingImpl()
	logging.Info("API開始")

	db := service.NewDynamoDBService()
	ocean, err := db.Fetch("伊豆海洋公園")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 503,
			Body:       err.Error(),
		}, nil
	}

	bytes, _ := json.Marshal(&ocean)
	logging.Info(string(bytes))

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bytes),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
