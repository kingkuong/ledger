#!/bin/sh
# Wait for DynamoDB to be ready
until curl -s http://dynamo:8000 > /dev/null; do
    echo "Waiting for DynamoDB..."
    sleep 2
done

echo "Creating DynamoDB tables..."
aws dynamodb create-table \
    --table-name transactions \
    --attribute-definitions \
        AttributeName=PK,AttributeType=S \
        AttributeName=SK,AttributeType=S \
    --key-schema \
        AttributeName=PK,KeyType=HASH \
        AttributeName=SK,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://dynamo:8000 \
    --region us-east-1

echo "Tables created successfully!"
