#ARG GO_VERSION=latest
FROM golang:latest AS builder
ARG SYSTEM=LongMeanReversion

WORKDIR /app
ADD go.mod go.sum ${SYSTEM} ./
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:latest
COPY --from=builder /app/main .

