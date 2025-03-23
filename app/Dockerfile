FROM golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o main .

FROM scratch

COPY --from=builder /app/main /main
COPY --from=builder /app/aws_env.yml /env.yml

EXPOSE 5050
ENTRYPOINT ["/main"]
