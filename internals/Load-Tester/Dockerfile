FROM golang:1.22-alpine

WORKDIR /workspace

COPY internals/Config Config/
COPY internals/Service Service/
COPY internals/Statistics Statistics/
# COPY ../bench.json /workspace/bench.json
COPY ../.env /workspace/.env

WORKDIR /workspace/Load-Tester

COPY internals/Load-Tester/go.mod internals/Load-Tester/go.sum ./

RUN go mod tidy
RUN go mod download

COPY internals/Load-Tester .  

RUN go build -o /workspace/Load-Tester/load-tester main.go

EXPOSE 8080
CMD ["./load-tester", "/workspace/bench.json"]