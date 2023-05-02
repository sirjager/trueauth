FROM golang:1.19.7-alpine3.17 as builder
WORKDIR /app

RUN apk update
RUN apk add libc6-compat git build-base

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd
COPY ./config ./config
COPY ./docs ./docs
COPY ./internal ./internal
COPY ./pkg ./pkg

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -o main ./cmd/main.go

FROM alpine:latest as runner
WORKDIR /app

RUN apk add libc6-compat

EXPOSE 4420
EXPOSE 4421

COPY  start.sh .
COPY  wait-for.sh .
COPY example.env .
COPY ./migrations ./migrations

RUN chmod +x start.sh wait-for.sh

COPY --from=builder /app/main .

ENTRYPOINT [ "/app/main" ]
