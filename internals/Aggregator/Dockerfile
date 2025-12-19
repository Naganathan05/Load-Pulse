FROM golang:1.22-alpine

WORKDIR /workspace

COPY internals/Config Config/
COPY internals/Service Service/
COPY internals/Statistics Statistics/
COPY ../testConfig.json /workspace/testConfig.json
COPY ../.env /workspace/.env

WORKDIR /workspace/Aggregator

COPY internals/Aggregator/go.mod internals/Aggregator/go.sum ./

RUN go mod tidy
RUN go mod download

COPY internals/Aggregator .  

RUN go build -o /workspace/Aggregator/aggregator main.go

EXPOSE 8081
CMD ["./aggregator", "/workspace/testConfig.json"]