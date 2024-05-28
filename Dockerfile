FROM golang:1.22-alpine as builder
WORKDIR /app

RUN apk update
RUN apk add libc6-compat git build-base

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -o main main.go

FROM alpine:latest as runner
WORKDIR /app

RUN apk add libc6-compat

EXPOSE 4420
EXPOSE 4421

COPY  start.sh .
COPY  wait-for.sh .
COPY defaults.env .

RUN chmod +x start.sh wait-for.sh

COPY --from=builder /app/main .

ENTRYPOINT [ "/app/main" ]

LABEL name="trueauth"
LABEL org.opencontainers.image.source="https://github.com/sirjager/trueauth"
