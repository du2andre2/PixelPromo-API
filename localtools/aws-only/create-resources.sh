#!/bin/bash

echo "Criando Tabelas..."

# Criando a tabela pp-user-catalog
aws dynamodb create-table \
    --table-name pp-user-catalog \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=createdAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes '[{
        "IndexName": "CreatedAtIndex",
        "KeySchema": [{"AttributeName": "createdAt", "KeyType": "HASH"}],
        "Projection": {"ProjectionType": "ALL"}
    }]' \
    --profile=admin --region=us-east-1 > /dev/null

# Criando a tabela pp-promotion-catalog
aws dynamodb create-table \
    --table-name pp-promotion-catalog \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=createdAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes '[{
        "IndexName": "CreatedAtIndex",
        "KeySchema": [{"AttributeName": "createdAt", "KeyType": "HASH"}],
        "Projection": {"ProjectionType": "ALL"}
    }]' \
    --profile=admin --region=us-east-1 > /dev/null

# Criando a tabela pp-promotion-interaction
aws dynamodb create-table \
    --table-name pp-promotion-interaction \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=createdAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes '[{
        "IndexName": "CreatedAtIndex",
        "KeySchema": [{"AttributeName": "createdAt", "KeyType": "HASH"}],
        "Projection": {"ProjectionType": "ALL"}
    }]' \
    --profile=admin --region=us-east-1 > /dev/null

# Criando a tabela pp-category-catalog
aws dynamodb create-table \
    --table-name pp-category-catalog \
    --attribute-definitions \
        AttributeName=name,AttributeType=S \
    --key-schema \
        AttributeName=name,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --profile=admin --region=us-east-1 > /dev/null

# Criando a tabela pp-user-score
aws dynamodb create-table \
    --table-name pp-user-score \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=createdAt,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes '[{
        "IndexName": "CreatedAtIndex",
        "KeySchema": [{"AttributeName": "createdAt", "KeyType": "HASH"}],
        "Projection": {"ProjectionType": "ALL"}
    }]' \
    --profile=admin --region=us-east-1 > /dev/null

echo "Criando Buckets..."

aws s3api create-bucket \
    --bucket pp-user-pictures \
    --profile=admin --region=us-east-1 > /dev/null

aws s3api create-bucket \
    --bucket pp-promotion-images \
    --profile=admin --region=us-east-1 > /dev/null

echo "Configuração finalizada com sucesso!"
