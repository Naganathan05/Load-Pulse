# Load-Pulse
Load testing tool built in Go which works based on the Raft consensus algorithm. Designed for scalable and distributed benchmarking using Dockerized microservices.

## High-Level Architecture Diagram
![image](https://github.com/user-attachments/assets/9643db11-e4df-48b2-b642-666e5168ff8f)

---

## 🚀 Setup Guide

### 🔧 Prerequisites
Before using Load Pulse, ensure the following tools are installed:
- **[Docker](https://www.docker.com/get-started)** – Containerization engine
- **[Git](https://git-scm.com/)** – Version control system

---

### 1. Clone the Repository
```sh
git clone https://github.com/Naganathan05/Load-Pulse.git
cd Load-Pulse
```

---

### 2. Modify `bench.json`
- Edit the `bench.json` file to configure the load test parameters and target server.
- If the target server runs on **localhost**, set the host like this:
```json
"host": "http://host.docker.internal:<port_number>/"
```
> This allows containers to reach your local server properly.

---

### 3. Start Docker Desktop
Ensure Docker Desktop is up and running, as the `docker-compose` command depends on the Docker daemon.

---

### 4. Start the Load Test
Start the tool using:
```sh
go run .\main.go run
```
This command spins up all required microservice containers (load tester, aggregator, Redis, RabbitMQ) and initiates the benchmarking process.

---

### 5. View Load Testing Results
- Once the test completes, results will be logged by the aggregator container.
- The following metrics are recorded and printed:
  - **Average Response Time**
  - **Max Response Time**
  - **Min Response Time**
  - **Total Requests Sent**
  - **Successful vs Failed Requests**

---

## 🛠 Troubleshooting

- **Containers fail to start?**
  - Run `docker ps -a` to check logs and exit codes.

- **Unable to reach localhost APIs?**
  - Always use `host.docker.internal:<port_number>` inside container configs.

---

## Maintainer
`Naganathan M R` 
