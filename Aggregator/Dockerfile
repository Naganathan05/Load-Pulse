FROM golang:1.22-alpine

WORKDIR /workspace

COPY ../Config Config/
COPY ../Service Service/
COPY ../Statistics Statistics/
COPY bench.json /workspace/bench.json
COPY .env /workspace/.env

WORKDIR /workspace/Aggregator

COPY Aggregator/go.mod Aggregator/go.sum ./

RUN go mod tidy
RUN go mod download

COPY Aggregator .  

RUN go build -o /workspace/Aggregator/aggregator main.go

EXPOSE 8081
CMD ["./aggregator", "/workspace/bench.json"]