FROM golang:1.22.1-alpine3.18

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

WORKDIR /app/cmd

RUN go build -o websocket-server

EXPOSE 8000

CMD ["./websocket-server"] 6