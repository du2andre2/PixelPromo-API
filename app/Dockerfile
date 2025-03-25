FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main /main
COPY --from=builder /app/env.yml /env.yml
COPY --from=builder /app/configmap.yml /cm.yml

EXPOSE 5050
ENTRYPOINT ["/main"]
