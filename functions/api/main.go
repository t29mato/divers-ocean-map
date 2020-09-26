package main

import (
	"api/logging"
	"api/model"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler ...
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logging := logging.NewOceanLoggingImpl()
	logging.Info("API開始")

	ocean := model.NewOcean("適当な地名")
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
