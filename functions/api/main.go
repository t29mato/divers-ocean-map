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

	resource := request.Resource
	logging.Info("API開始, path=", resource)

	db := service.NewDynamoDBDatabase(logging)

	switch resource {
	case "/api/oceans/{name}":
		locationName := request.PathParameters["name"]
		ocean, err := db.FetchLatestOcean(locationName)
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
	case "/api/oceans":
		oceans, err := db.FetchAllLatestOceans()
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}

		bytes, _ := json.Marshal(&oceans)
		logging.Info(string(bytes))

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(bytes),
		}, nil

	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotAcceptable,
		Body:       "無効なURLです",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
