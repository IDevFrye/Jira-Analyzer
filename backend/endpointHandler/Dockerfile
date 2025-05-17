# Dockerfile (backend)
FROM golang:1.24.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend ./

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=build /app/backend .
COPY wait-for-it.sh .             
RUN chmod +x wait-for-it.sh       

CMD ["./wait-for-it.sh", "postgres:5432", "--", "./backend"]