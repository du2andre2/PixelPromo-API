
aws dynamodb put-item \
    --table-name pp-category-catalog \
    --item \
        '{"name": {"S": "rpg"}}'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-category-catalog \
    --item \
        '{"name": {"S": "fps"}}'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-category-catalog \
    --item \
        '{"name": {"S": "estrategia"}}'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-category-catalog \
    --item \
        '{"name": {"S": "indie"}}'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-category-catalog \
    --item \
        '{"name": {"S": "plataforma"}}'\
    --endpoint-url http://localhost:4566 > /dev/null