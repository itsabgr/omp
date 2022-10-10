# syntax=docker/dockerfile:1

FROM golang:1.19-alpine
WORKDIR /app
COPY . .
EXPOSE 4444
CMD go run "./cmd/service"