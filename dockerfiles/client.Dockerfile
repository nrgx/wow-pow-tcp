# STAGE BUILD
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -o client ./cmd/client

# STAGE MAIN
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/client /app/client

CMD [ "/app/client" ]
