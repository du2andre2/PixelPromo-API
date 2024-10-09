
docker-compose down
docker-compose up -d

sleep 5

echo Criando Tabelas

aws dynamodb create-table \
    --table-name pp-user-catalog \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=CreatedAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
        AttributeName=DT,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:4566 > /dev/null

aws s3api create-bucket \
    --bucket pp-user-pictures \
    --endpoint-url http://localhost:4566 > /dev/null


aws dynamodb create-table \
    --table-name pp-promotion-catalog \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=CreatedAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
        AttributeName=DT,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb create-table \
    --table-name pp-promotion-interaction \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=CreatedAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
        AttributeName=DT,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb create-table \
    --table-name pp-category-catalog \
    --attribute-definitions \
        AttributeName=name,AttributeType=S \
    --key-schema \
        AttributeName=name,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb create-table \
    --table-name pp-user-score \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=CreatedAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
        AttributeName=DT,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:4566 > /dev/null


echo tabelas criadas:

aws dynamodb list-tables --endpoint-url http://localhost:4566

echo criando buckets

aws s3api create-bucket \
    --bucket pp-user-pictures \
    --endpoint-url http://localhost:4566 > /dev/null

aws s3api create-bucket \
    --bucket pp-promotion-images \
    --endpoint-url http://localhost:4566 > /dev/null
