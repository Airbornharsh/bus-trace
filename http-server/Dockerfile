FROM golang:1.22.1-alpine3.18

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

WORKDIR /app/cmd

RUN go build -o bus-trace-http-server-go

EXPOSE 8001

CMD ["./bus-trace-http-server-go"]