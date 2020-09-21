aws dynamodb create-table --table-name Oceans \
    --attribute-definitions \
        AttributeName=Name,AttributeType=S \
        AttributeName=MeasuredTime,AttributeType=S \
    --key-schema \
        AttributeName=Name,KeyType=HASH \
        AttributeName=MeasuredTime,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:8000