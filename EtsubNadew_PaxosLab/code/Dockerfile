FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o paxos-web web-service.go

EXPOSE 8080

CMD ["./paxos-web"]
