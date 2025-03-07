FROM golang:1.22.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM debian:latest

WORKDIR /app

RUN apt-get update && apt-get install -y curl netcat-openbsd

COPY --from=builder /app/main .

EXPOSE 8080

CMD bash -c 'until nc -z mysql 3306; do sleep 1; done; echo "MySQL is up!"; ./main'
