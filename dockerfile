FROM golang:1.23.1-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/main /main

EXPOSE 80

ENTRYPOINT ["./main"]