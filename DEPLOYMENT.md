# Sentinel Quick Start (Docker Image)

Sentinel is a lightweight server and container monitoring tool.

With a **single Docker image**, you can run:

* **Master** → Web dashboard + API
* **Agent** → Collects server/container metrics and sends to master

Because apparently one image doing two jobs is now the sensible timeline.

---

# Prerequisites

* Docker installed
* Linux/macOS/Windows with Docker Desktop
* Network access between Agent servers and Master server
* Docker socket access for Agent container stats

---

# Pull or Build Image

## Pull

```bash
docker pull sentinel:latest
```

## Or Build Locally

```bash
docker build -t sentinel .
```

---

# Architecture

```text
Agent(s) ---> NATS Broker ---> Master Dashboard
```

Recommended setup:

```text
Master Server:
  NATS
  Sentinel Master

Remote Servers:
  Sentinel Agent
```

---

# Step 1: Start NATS Broker

Run on the same server as Master (recommended):

```bash
docker run -d \
  --name nats \
  -p 4222:4222 \
  nats
```

---

# Step 2: Start Sentinel Master

```bash
docker run -d \
  --name sentinel-master \
  -p 8080:8080 \
  -e NATS_URL=nats://host.docker.internal:4222 \
  sentinel master
```

## Open Dashboard

```text
http://localhost:8080
```

---

# Step 3: Start Sentinel Agent

Run this on any server you want to monitor.

```bash
docker run -d \
  --name sentinel-agent \
  -e SERVER_ID=server-1 \
  -e NATS_URL=nats://MASTER_SERVER_IP:4222 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  sentinel agent
```

---

# Example Multi-Server Setup

## Master Server IP

```text
10.72.20.38
```

## Agent Command

```bash
docker run -d \
  --name sentinel-agent \
  -e SERVER_ID=prod-node-1 \
  -e NATS_URL=nats://10.72.20.38:4222 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  sentinel agent
```

---

# Command Modes

## Run Master

```bash
docker run sentinel master
```

## Run Agent

```bash
docker run sentinel agent
```

---

# Environment Variables

## Common

| Variable   | Description         | Example                   |
| ---------- | ------------------- | ------------------------- |
| `NATS_URL` | NATS broker address | `nats://10.72.20.38:4222` |

## Agent Only

| Variable    | Description           | Example       |
| ----------- | --------------------- | ------------- |
| `SERVER_ID` | Unique server name/id | `prod-node-1` |

---

# Docker Socket Mount (Agent)

Required for container stats:

```bash
-v /var/run/docker.sock:/var/run/docker.sock
```

Without it, container metrics will be fiction.

---

# Logs

## Master Logs

```bash
docker logs -f sentinel-master
```

## Agent Logs

```bash
docker logs -f sentinel-agent
```

---

# Stop Containers

```bash
docker stop sentinel-master sentinel-agent nats
docker rm sentinel-master sentinel-agent nats
```

---

# Troubleshooting

## NATS Connection Error

```text
nats: no servers available for connection
```

Check:

* Correct `NATS_URL`
* Port `4222` open
* NATS container running

---

## No Container Data

Ensure Docker socket mounted:

```bash
-v /var/run/docker.sock:/var/run/docker.sock
```

---

## Agent Shows Offline

Check:

* Agent container running
* Network connectivity to NATS
* Server clock/time sync

---

# Recommended Production Setup

* Run Master + NATS on dedicated server
* Use private IP / VPN
* Use unique `SERVER_ID` for each agent
* Enable restart policy:

```bash
--restart unless-stopped
```

---

# Example Production Agent

```bash
docker run -d \
  --restart unless-stopped \
  --name sentinel-agent \
  -e SERVER_ID=prod-api-01 \
  -e NATS_URL=nats://10.72.20.38:4222 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  sentinel agent
```

---

# Summary

## Start Master Side

```bash
docker run -d --name nats -p 4222:4222 nats
docker run -d --name sentinel-master -p 8080:8080 -e NATS_URL=nats://host.docker.internal:4222 sentinel master
```

## Start Agent Side

```bash
docker run -d --name sentinel-agent -e SERVER_ID=node-1 -e NATS_URL=nats://MASTER_IP:4222 -v /var/run/docker.sock:/var/run/docker.sock sentinel agent
```

---

