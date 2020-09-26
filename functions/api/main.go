package main

import (
	"api/logging"
	"api/service"
	"encoding/json"
	"net/http"

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
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}

		bytes, _ := json.Marshal(&ocean)
		logging.Info(string(bytes))

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(bytes),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotAcceptable,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
