FROM golang:1.22-alpine

WORKDIR /workspace

COPY ../Config Config/
COPY ../Service Service/
COPY ../Statistics Statistics/
COPY bench.json /workspace/bench.json
COPY .env /workspace/.env

WORKDIR /workspace/Load-Tester

COPY Load-Tester/go.mod Load-Tester/go.sum ./

RUN go mod tidy
RUN go mod download

COPY Load-Tester .  

RUN go build -o /workspace/Load-Tester/load-tester main.go

EXPOSE 8080
CMD ["./load-tester", "/workspace/bench.json"]