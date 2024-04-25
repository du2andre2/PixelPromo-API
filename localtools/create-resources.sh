
docker-compose down
docker-compose up -d

echo Criando Tabelas

aws dynamodb create-table \
--table-name Music \
--attribute-definitions \
AttributeName=Artist,AttributeType=S \
AttributeName=SongTitle,AttributeType=S \
--key-schema \
AttributeName=Artist,KeyType=HASH \
AttributeName=SongTitle,KeyType=RANGE \
--provisioned-throughput \
ReadCapacityUnits=5,WriteCapacityUnits=5 \
--stream-specification \
StreamEnabled=false \
--endpoint-url http://localhost:4566


echo tabelas criadas

aws dynamodb list-tables --endpoint-url http://localhost:4566

echo criando itens teste

aws dynamodb put-item \
    --table-name Music  \
    --item \
        '{"Artist": {"S": "No One You Know"}, "SongTitle": {"S": "Call Me Today"}, "AlbumTitle": {"S": "Somewhat Famous"}, "Awards": {"N": "1"}}' \
    --endpoint-url http://localhost:4566