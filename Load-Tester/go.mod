module loadpulse.local/Load-Tester

go 1.22.5

require (
	loadpulse.local/Config v0.0.0
	loadpulse.local/Service v0.0.0
	loadpulse.local/Statistics v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
)

replace loadpulse.local/Aggregator => ../Aggregator

replace loadpulse.local/Service => ../Service

replace loadpulse.local/Config => ../Config

replace loadpulse.local/Statistics => ../Statistics
