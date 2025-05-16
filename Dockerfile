FROM golang:latest

RUN apt-get update && apt-get install -y postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o main .

WORKDIR /app
COPY wait-for.sh /wait-for.sh
RUN chmod +x /wait-for.sh

CMD ["/wait-for.sh", "postgres", "5432", "./cmd/main"]