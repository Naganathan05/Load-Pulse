# Load-Pulse
Load Testing tool built in Go which works based on Raft Algorithm.


## High Level Architecture Diagram:
![image](https://github.com/user-attachments/assets/9643db11-e4df-48b2-b642-666e5168ff8f)

# Setup

## Prerequisites
Before using Load Pulse, ensure you have the following installed:
- **Docker** (Install from [Docker Official Website](https://www.docker.com/get-started))
- **Git** (Install from [Git Official Website](https://git-scm.com/))

## Setup Guide

### 1. Clone the Repository
```sh
git clone https://github.com/Naganathan05/Load-Pulse.git
cd Load-Pulse
```

### 2. Create the `.env` File
In the root directory of the project, create a file named `.env` and add the following content:

```ini
REDIS_KEY = "concurrencyCount"
BASE_QUEUE_NAME = "StatsEventQueue"
CLUSTER_SIZE = 10  # Maximum Allowed Workers Per Cluster
REQUEST_SLEEP_TIME = 50

# Service URLs
REDIS_URL_LOCAL = "localhost:6379"
REDIS_URL_DOCKER = "redis:6379"
RABBITMQ_URL_LOCAL = "amqp://guest:guest@localhost:5672/"
RABBITMQ_URL_DOCKER = "amqp://guest:guest@rabbitmq:5672/"

# Ports
LOAD_TESTER_PORT = "8080"
AGGREGATOR_PORT = "8081"
```

### 3. Modify `bench.json`
- Edit the `bench.json` file to configure the target endpoint and server.
- If the target server is running on **localhost**, set the host as:
  ```json
  "host": "http://host.docker.internal:<port_number>/"
  ```
  This ensures proper communication between the container and the host machine.

### 4. Start the Containers
Run the following command to start the services:
```sh
docker compose up -d
```
This will start all necessary containers in detached mode.

### 5. Check Aggregator Logs
To verify that the load testing tool is running correctly, check the logs of the **aggregator** container:
```sh
docker logs aggregator
```

### 6. View Load Testing Results
- Once the load test runs, check the results for the respective endpoints.
- The statistics of the tested endpoints will be available in the logs and database.

## Troubleshooting
- **Containers fail to start?** Run `docker ps -a` to check for errors.
- **RabbitMQ or Redis not connecting?** Ensure the correct service URL is being used from the `.env` file.
- **Unable to reach localhost APIs?** Use `host.docker.internal:<port_number>` inside the container.

---
Maintainer: `Naganathan M R`
