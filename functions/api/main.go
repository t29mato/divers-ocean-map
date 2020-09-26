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

	resource := request.Resource

	db := service.NewDynamoDBService()

	switch resource {
	case "/api/oceans/{name}":
		locationName := request.PathParameters["name"]
		ocean, err := db.Fetch(locationName)
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
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "hoge",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
