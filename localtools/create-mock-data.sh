echo "Creating mock data..."

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

aws dynamodb get-tables \
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


aws s3 cp ./imgs s3://pp-user-pictures/ --recursive \
    --endpoint-url http://localhost:4566


aws dynamodb put-item \
    --table-name pp-user-catalog \
    --item \
        '{
            "id": {"S":"1"},
            "email": {"S":"edu@gmail.com"},
            "name": {"S":"edu"},
            "password": {"S":"12345678"},
            "pictureUrl": {"S":"http://localhost:4566/pp-user-pictures/perfil1.jpg"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-user-catalog \
    --item \
        '{
            "id": {"S":"2"},
            "email": {"S":"lucas@gmail.com"},
            "name": {"S":"lucas"},
            "password": {"S":"12345678"},
            "pictureUrl": {"S":"http://localhost:4566/pp-user-pictures/perfil2.jpg"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-user-catalog \
    --item \
        '{
            "id": {"S":"3"},
            "email": {"S":"joao@gmail.com"},
            "name": {"S":"joao"},
            "password": {"S":"12345678"},
            "pictureUrl": {"S":"http://localhost:4566/pp-user-pictures/perfil3.jpg"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null


aws dynamodb put-item \
    --table-name pp-user-catalog \
    --item \
        '{
            "id": {"S":"4"},
            "email": {"S":"pedro@gmail.com"},
            "name": {"S":"pedro"},
            "password": {"S":"12345678"},
            "pictureUrl": {"S":"http://localhost:4566/pp-user-pictures/perfil4.jpg"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null


aws s3 cp ./imgs s3://pp-promotion-images/ --recursive \
    --endpoint-url http://localhost:4566


aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"1"},
            "userId": {"S":"1"},
            "title": {"S":"jogo 1"},
            "description": {"S":"jogo 1 na promoção"},
            "categories": {"L": [ {"S": "fps"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo1.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"300,00"},
            "discountedPrice": {"S":"255,00"},
            "discountBadge": {"S":"15"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"2"},
            "userId": {"S":"1"},
            "title": {"S":"jogo 2"},
            "description": {"S":"jogo 2 na promoção"},
            "categories": {"L": [ {"S": "estrategia"} , {"S": "indie"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo2.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"100,00"},
            "discountedPrice": {"S":"91,00"},
            "discountBadge": {"S":"9"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"3"},
            "userId": {"S":"1"},
            "title": {"S":"jogo 3"},
            "description": {"S":"jogo 3 na promoção"},
            "categories": {"L": [ {"S": "plataforma"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo3.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"100,00"},
            "discountedPrice": {"S":"86,00"},
            "discountBadge": {"S":"14"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"4"},
            "userId": {"S":"2"},
            "title": {"S":"jogo 4"},
            "description": {"S":"jogo 4 na promoção"},
            "categories": {"L": [ {"S": "fps"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo1.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"300,00"},
            "discountedPrice": {"S":"150,00"},
            "discountBadge": {"S":"50"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"5"},
            "userId": {"S":"2"},
            "title": {"S":"jogo 5"},
            "description": {"S":"jogo 5 na promoção"},
            "categories": {"L": [ {"S": "estrategia"} , {"S": "indie"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo2.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"300,00"},
            "discountedPrice": {"S":"200,00"},
            "discountBadge": {"S":"33"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"6"},
            "userId": {"S":"2"},
            "title": {"S":"jogo 6"},
            "description": {"S":"jogo 6 na promoção"},
            "categories": {"L": [ {"S": "plataforma"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo3.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"100,00"},
            "discountedPrice": {"S":"75,00"},
            "discountBadge": {"S":"25"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"7"},
            "userId": {"S":"3"},
            "title": {"S":"jogo 7"},
            "description": {"S":"jogo 7 na promoção"},
            "categories": {"L": [ {"S": "fps"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo1.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"250,00"},
            "discountedPrice": {"S":"187,5"},
            "discountBadge": {"S":"25"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"8"},
            "userId": {"S":"3"},
            "title": {"S":"jogo 8"},
            "description": {"S":"jogo 8 na promoção"},
            "categories": {"L": [ {"S": "estrategia"} , {"S": "indie"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo2.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"100,00"},
            "discountedPrice": {"S":"50,00"},
            "discountBadge": {"S":"50,00"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"9"},
            "userId": {"S":"3"},
            "title": {"S":"jogo 9"},
            "description": {"S":"jogo 9 na promoção"},
            "categories": {"L": [ {"S": "plataforma"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo3.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"150,00"},
            "discountedPrice": {"S":"50,00"},
            "discountBadge": {"S":"66"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"10"},
            "userId": {"S":"4"},
            "title": {"S":"jogo 10"},
            "description": {"S":"jogo 10 na promoção"},
            "categories": {"L": [ {"S": "fps"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo1.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"150,00"},
            "discountedPrice": {"S":"100,00"},
            "discountBadge": {"S":"33"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"11"},
            "userId": {"S":"4"},
            "title": {"S":"jogo 11"},
            "description": {"S":"jogo 11 na promoção"},
            "categories": {"L": [ {"S": "estrategia"} , {"S": "indie"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo2.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"200,00"},
            "discountedPrice": {"S":"100,00"},
            "discountBadge": {"S":"50"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

aws dynamodb put-item \
    --table-name pp-promotion-catalog \
    --item \
        '{
            "id": {"S":"12"},
            "userId": {"S":"4"},
            "title": {"S":"jogo 12"},
            "description": {"S":"jogo 12 na promoção"},
            "categories": {"L": [ {"S": "plataforma"} , {"S": "rpg"}] },
            "imageUrl": {"S":"http://localhost:4566/pp-promotion-images/jogo3.jpg"},
            "link": {"S":"https://www.google.com"},
            "originalPrice": {"S":"100,00"},
            "discountedPrice": {"S":"60,00"},
            "discountBadge": {"S":"40"},
            "platform": {"S":"Steam"},
            "createdAt": {"S":"2024-09-01T12:34:32.5657674-03:00"}
        }'\
    --endpoint-url http://localhost:4566 > /dev/null

