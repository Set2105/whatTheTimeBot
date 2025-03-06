FROM golang:1.23.5-alpine3.21 AS builder

WORKDIR /app

COPY . .
ENV GOOS=linux GOARCH=amd64
RUN go mod tidy
RUN go build -ldflags="-s -w" -o bot ./cmd/main.go 

FROM alpine:3.21

WORKDIR /app
COPY --from=builder /app/bot .

CMD [ "./bot" ]