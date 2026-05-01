# Sentinel 

Sentinel is a lightweight, distributed system and Docker monitoring tool. It follows a master-agent architecture, using **NATS** as a high-performance messaging backbone and a **Preact**-based frontend for a responsive, real-time web dashboard.

## Features

- **Distributed Monitoring**: Deploy agents on multiple servers to collect metrics centrally.
- **System Metrics**: Real-time tracking of CPU, Memory, Disk usage, and IPv4 addresses.
- **Docker Integration**: Monitor container status, resource utilization, and perform actions like fetching logs or inspecting containers directly from the UI.
- **Live Dashboard**: A sleek, real-time web interface built with Go templates and Preact, featuring advanced sorting, search, and interactive container management modals.
- **High Performance**: Built with Go and NATS for low latency and minimal resource footprint.

## Architecture

1.  **NATS Broker**: The communication hub where agents publish data and the master subscribes.
2.  **Sentinel Agent**: Runs on target hosts. Collects system and Docker metrics and publishes them to NATS.
3.  **Sentinel Master**: Central server that collects data from NATS, stores it in-memory, and serves the Web UI/API.

## Tech Stack

- **Backend**: [Go](https://go.dev/)
- **Messaging**: [NATS](https://nats.io/)
- **Frontend**: [Preact](https://preactjs.com/), Go Templates, Vanilla CSS
- **Metrics**: `gopsutil`, Docker SDK

## 🏁 Getting Started

### 0. Build Images
Before running, build the Docker images for the master and agent:
```bash
# Build Master
docker build --target master -t sentinel-master .

# Build Agent
docker build --target agent -t sentinel-agent .
```

### 1. Simple Start (All Components)
The easiest way is using the provided Docker Compose file which starts NATS, the Sentinel Master, and a Sentinel Agent:
```bash
docker-compose up -d
```
The dashboard will be available at `http://localhost:8080`.

### 2. Distributed Deployment (Different Machines)

To monitor multiple servers, you need a central NATS broker accessible by all machines.

#### On the Central Server (Master + NATS):
1. Start NATS and Master using Docker:
   ```bash
   # In docker-compose.yaml, ensure nats port 4222 is exposed
   docker-compose up -d nats master
   ```
2. Note the public IP of this server (e.g., `192.168.1.50`).

#### On Remote Servers (Agents):
You can run only the agent using Docker:
```bash
docker run -d \
  --name sentinel-agent \
  --hostname remote-server-01 \
  -e NATS_URL="nats://192.168.1.50:4222" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --restart unless-stopped \
  <your-docker-registry>/sentinel-agent:latest
```
*(Note: You will need to build and push the agent image to a registry or copy the image to the remote machine.)*

## Configuration

Both components can be configured via environment variables:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `NATS_URL` | URL of the NATS server | `nats://localhost:4222` |
| `ADDR` | (Master only) Address to bind the HTTP server | `:8080` |

## API Endpoints

Sentinel Master provides a simple JSON API for integration with other tools:

- `GET /api/servers`: Lists all active servers and their latest system metrics.
- `GET /api/containers?ServerID=<id>`: Lists Docker containers and stats for a specific server.
- `POST /api/actions`: Submit an action (e.g., logs, inspect) to be executed by an agent on a specific container.
- `GET /api/result?hostname=<id>&containerId=<id>&action=<action>`: Poll for the result of a submitted action.

---
Built with ❤️ using Go and Preact.
