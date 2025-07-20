FROM golang:1.23.11-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["go", "run", "./src/cmd/main.go"]

