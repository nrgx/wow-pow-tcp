# STAGE BUILD
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -o server ./cmd/server

# STAGE MAIN
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server /app/server

EXPOSE 9090

CMD [ "/app/server" ]
