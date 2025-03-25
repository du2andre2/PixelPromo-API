#!/bin/bash

echo "Creating mock data..."

# Configurações
DYNAMODB_ENDPOINT="http://localhost:4566"
S3_BUCKET_USER="s3://pp-user-imgs"
S3_BUCKET_PROMOTION="s3://pp-promotion-imgs"
CATEGORY_COUNT=10
USER_COUNT=10
PROMOTION_COUNT=30
IMAGES_PATH="./imgs/promotions"
PICTURES_PATH="./imgs/users"

# Criar categorias
echo "Creating categories..."
for i in $(seq 1 $CATEGORY_COUNT); do
    CATEGORY_NAME="categoria_$i"
    aws dynamodb put-item \
        --table-name pp-category-catalog \
        --item "{\"name\": {\"S\": \"$CATEGORY_NAME\"}}" \
        --endpoint-url $DYNAMODB_ENDPOINT > /dev/null
done

# Enviar imagens para S3
echo "Uploading images to S3..."
aws s3 cp $PICTURES_PATH $S3_BUCKET_USER --recursive --endpoint-url $DYNAMODB_ENDPOINT > /dev/null
aws s3 cp $IMAGES_PATH $S3_BUCKET_PROMOTION --recursive --endpoint-url $DYNAMODB_ENDPOINT > /dev/null

# Criar usuários
echo "Creating users..."
for i in $(seq 1 $USER_COUNT); do
    USER_EMAIL="user$i@gmail.com"
    USER_NAME="user_$i"
    USER_PASSWORD="123123"
    USER_PICTURE="http://localhost:4566/pp-user-pictures/perfil$i.png"
    CREATED_AT=$(date -Iseconds)

    aws dynamodb put-item \
        --table-name pp-user-catalog \
        --item \
        "{
            \"id\": {\"S\":\"$i\"},
            \"email\": {\"S\":\"$USER_EMAIL\"},
            \"name\": {\"S\":\"$USER_NAME\"},
            \"password\": {\"S\":\"$USER_PASSWORD\"},
            \"pictureUrl\": {\"S\":\"$USER_PICTURE\"},
            \"createdAt\": {\"S\":\"$CREATED_AT\"}
        }" \
        --endpoint-url $DYNAMODB_ENDPOINT > /dev/null
done

# Criar promoções
echo "Creating promotions..."
for i in $(seq 1 $PROMOTION_COUNT); do
    PROMO_ID=$i
    USER_ID=$(( (i % USER_COUNT) + 1 ))  # Relacionar promoções com usuários existentes
    TITLE="Promoção Jogo $PROMO_ID"
    DESCRIPTION="Descrição da promoção $PROMO_ID"
    IMAGE_URL="http://localhost:4566/pp-promotion-images/jogo$(( (i % 12) + 1 )).png"
    LINK="https://example.com/promo_$PROMO_ID"

    # Gerar preços em float
    ORIGINAL_PRICE=$(awk -v min=50 -v range=500 'BEGIN { printf "%.2f", min + rand() * range }')
    DISCOUNT_PERCENT=$(awk 'BEGIN { printf "%.2f", (10 + rand() * 40) }')  # Desconto de 10% a 50%
    DISCOUNTED_PRICE=$(awk -v op="$ORIGINAL_PRICE" -v dp="$DISCOUNT_PERCENT" 'BEGIN { printf "%.2f", op * (1 - dp / 100) }')

    # Cálculo do badge de desconto
    DISCOUNT_BADGE=$(awk -v dp="$DISCOUNT_PERCENT" 'BEGIN { printf "%.0f", dp }')

    PLATFORM="Steam"
    CREATED_AT=$(date -Iseconds)
    CATEGORY_COUNT=$((RANDOM % 3 + 1))
    CATEGORIES=""

    for j in $(seq 1 $CATEGORY_COUNT); do
        CATEGORY_INDEX=$(( (i + j) % CATEGORY_COUNT + 1 ))
        CATEGORIES="${CATEGORIES}{\"S\": \"categoria_$CATEGORY_INDEX\"},"
    done

    # Remover vírgula extra
    CATEGORIES=${CATEGORIES%,}

    aws dynamodb put-item \
        --table-name pp-promotion-catalog \
        --item \
        "{
            \"id\": {\"S\":\"$PROMO_ID\"},
            \"userId\": {\"S\":\"$USER_ID\"},
            \"title\": {\"S\":\"$TITLE\"},
            \"description\": {\"S\":\"$DESCRIPTION\"},
            \"categories\": {\"L\": [ $CATEGORIES ]},
            \"imageUrl\": {\"S\":\"$IMAGE_URL\"},
            \"link\": {\"S\":\"$LINK\"},
            \"originalPrice\": {\"N\":\"$ORIGINAL_PRICE\"},
            \"discountedPrice\": {\"N\":\"$DISCOUNTED_PRICE\"},
            \"discountBadge\": {\"N\":\"$DISCOUNT_BADGE\"},
            \"platform\": {\"S\":\"$PLATFORM\"},
            \"createdAt\": {\"S\":\"$CREATED_AT\"}
        }" \
        --endpoint-url $DYNAMODB_ENDPOINT > /dev/null
done

echo "Mock data creation complete!"
