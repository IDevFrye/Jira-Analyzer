FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o jiraConnector ./cmd/service

EXPOSE 8082

CMD ["./jiraConnector"]