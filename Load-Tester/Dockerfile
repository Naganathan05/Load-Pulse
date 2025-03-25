FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o load-tester.exe main.go

EXPOSE 8080

CMD ["./load-tester", "..\bench.json"]