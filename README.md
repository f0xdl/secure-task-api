
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![License](https://img.shields.io/github/license/Ileriayo/markdown-badges?style=for-the-badge)

# Practice: Secure Task API
## 📚Table of contents
- [Practice: Secure Task API](#practice-secure-task-api)
    * [📚Table of contents](#table-of-contents)
    * [📝Description](#description)
        + [Features](#features)
    * [🚀Getting Started](#getting-started)
        + [Prerequisites](#prerequisites)
        + [Clone](#clone)
        + [Environment](#environment)
        + [Launch with Docker](#launch-with-docker)
        + [Custom Launch](#custom-launch)
    * [Routes](#routes)
        + [Mock запрос на выполнения операций (without auth).](#mock-request-for-operations-without-auth)
        + [Mock запрос с авторизацией для получения метрик.](#mock-request-with-auth-to-get-metrics)
    * [🧪Testing](#testing)
    * [📁Folder structure](#folder-structure)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

## 📝Description
Secure Task API is a minimalistic HTTP server on Go designed to practice secure API design.
It includes a basic task system with 
administrative interface protection, 
logging, 
limiting of requests 
and deploy with Docker Compose.

The project is focused on practicing backend development in Golang:
- working with middleware (recovery, rate limiting, logging)
- HTTP API security (Basic Auth)
- lifecycle organization (graceful shutdown)
- preparation for deployment and monitoring (healthcheck, metrics)

### Features
- ❤️ Health and readiness endpoints
- 🔐 Auth-protected admin API (HTTP Basic)
- 📄 Console logging for HTTP requests
- 🧰 Configurable rate limiting (Redis or Memory)
- 🛡️ Recovery middleware
- 🧪 Test coverage for core logic
- 🧹 Graceful shutdown support
- 📦 Docker Compose setup

## 🚀Getting Started
### Prerequisites
- [Docker](https://www.docker.com/get-started/)
- [go 1.23](https://go.dev/dl/)
### Clone
To build the application, you need to download the current version of the repository:
```shell
git clone https://github.com/f0xdl/secure-task-api
cd secure-task-api
```
### Environment
Environment variables can be passed manually or set in the `.env` file.
```dotenv
# Server configuration
API_PREFIX=/api/v1
HOST=":8080"

# Cache
REDIS_HOST=redis:6379
REDIS_PASSWORD=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81

# Auth
AUTH_USERNAME=admin
AUTH_PASSWORD=T3st
```

> [!TIP]
> `.env.example` is located at the directory level.
> For a quick one, rename it to `.env`.

### Launch with Docker
At the root is the `docker-compose.yml` file that raises the local Redis server and the HTTP server
from `Dockerfile`. The build is done in two-stage mode to reduce the size of the final image.
For a quick build, you can use the command:
```shell
docker compose up -D
```
To stop services you can use the command:
```shell
docker compose down
```

### Custom Launch
1. It is necessary to raise Redis server
2. Load environment variables and start the application
```shell
source .env && go run  ./cmd/main.go
```
## Routes
### Mock request for operations (without auth)
```http request
POST http://localhost:8080/api/v1/task HTTP/1.1
Content-Type: application/json

{ "value": 32 }
```
- 422 Response `wrong JSON format`
- 200 Response
```json
{"result":1024}
```

### Mock request with auth to get metrics
```http request
GET http://localhost:8080/api/v1/admin/metrics HTTP/1.1
Content-Type: application/json
Authorization: <Base64 encoded username and password>
```
- 200 Response: `<empty response>` 
- 401 Response: `Unauthorized`
## 🧪Testing
Run unit tests with:
```shell
go test ./...
```

## 📁Folder structure
```shell
cmd/            # Entrypoint
internal/
  handlers/     
  middleware/   
  httpserver/
```