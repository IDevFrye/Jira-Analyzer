FROM golang:1.24.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend ./cmd/service

FROM debian:bullseye-slim
WORKDIR /app

COPY --from=build /app/backend .
COPY --from=build /app/wait-for-it.sh .

COPY --from=build /app/configs ./configs

RUN chmod +x wait-for-it.sh

CMD ["./wait-for-it.sh", "postgres:5432", "--", "./backend", "--config=./configs/config.yaml"]
