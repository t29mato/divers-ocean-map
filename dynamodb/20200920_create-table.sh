aws dynamodb create-table --table-name oceans \
    --attribute-definitions \
        AttributeName=name,AttributeType=S \
        AttributeName=measured_time,AttributeType=S \
    --key-schema \
        AttributeName=name,KeyType=HASH \
        AttributeName=measured_time,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:8000