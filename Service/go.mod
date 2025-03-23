module loadpulse.local/Service

go 1.22.5

require (
	github.com/redis/go-redis/v9 v9.7.3
	github.com/streadway/amqp v1.1.0
	loadpulse.local/Config v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace loadpulse.local/Config => ../Config

replace loadpulse.local/Statistics => ../Statistics
