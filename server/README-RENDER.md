# Deploying Barista Microservices to Render

This guide provides step-by-step instructions for deploying the Barista project as individual services on [Render](https://render.com). The architecture separates **Public Web Services** (accessible from the internet) from **Private Services** (internal gRPC/AMQP services).

## 1. Prerequisites
- A Render account.
- A GitHub repository with the Barista code.
- A PostgreSQL database (Render managed or external).
- A Redis instance (Render managed).
- A RabbitMQ instance (See [Section 4](#4-custom-rabbitmq-service) for custom deployment).

---

## 2. Public Web Services
These services will have a public URL (e.g., `https://api-gateway-xxx.onrender.com`).

### api-gateway
1. Click **New +** > **Web Service**.
2. Connect your repository.
3. **Region**: Choose the same region for all services.
4. **Runtime**: `Go`
5. **Root Directory**: `server/api-gateway`
6. **Build Command**: `go build -o server ./cmd/main.go`
7. **Start Command**: `./server`
8. **Environment Variables**:
   - `PORT`: `8080` (Render will map traffic to this)
   - `JWT_KEY`: A long random string.
   - `ALLOWED_ORIGINS`: Your frontend URL (e.g., `https://barista.onrender.com`). For multiple origins, use a comma-separated list: `https://app1.com, https://app2.com`.
   - `AUTH_SERVICE_URL`, `SHOP_SERVICE_URL`, etc.: Use the **Internal Service Name** (e.g., `http://auth-service:8080`).

### ws (WebSocket Service)
1. Click **New +** > **Web Service**.
2. **Runtime**: `Go`
3. **Root Directory**: `server/services/ws`
4. **Build Command**: `go build -o server ./cmd/main.go`
5. **Start Command**: `./server`
6. **Environment Variables**:
   - `PORT`: `8080`
   - `ALLOWED_ORIGINS`: Your frontend URL.

---

## 3. Private Services
These services are **not** accessible from the internet. They are reachable only via internal hostnames provided by Render.

### General Steps for Internal Services
(`auth`, `order`, `payment`, `session`, `shop`, `table`, `history`, `menu`, `qrgen`)

1. Click **New +** > **Private Service**.
2. Connect your repository.
3. **Runtime**: `Go`
4. **Root Directory**: `server/services/<service_name>`
5. **Build Command**: `go build -o server ./cmd/main.go`
6. **Start Command**: `./server`
7. **Advanced Settings**:
   - **Health Check Path**: `/health`
8. **Environment Variables**:
   - `PORT`: `8080` (Internal port)
   - `DB_URL`: PostgreSQL connection string.
   - `REDIS_URL`: (For `session` service) External/Internal Redis string.
   - `*_SERVICE_URL`: Internal URLs for other dependencies (e.g., `http://menu-service:8080`).

> [!TIP]
> **Internal Hostnames**: When service `A` needs to call service `B`, use `http://service-b-name:8080`. You can find the internal name in the service's dashboard.

---

## 4. Custom RabbitMQ Service
Deploy the specialized RabbitMQ service used by `table` and `qrgen`.

1. Click **New +** > **Private Service**.
2. Connect your repository.
3. **Name**: `rabbitmq-service` (or any name you prefer)
4. **Environment**: `Docker`
5. **Root Directory**: `server/services/rabbitmq`
6. **Dockerfile Path**: `./app.dockerfile`
7. **Advanced Settings**:
   - **Health Check Path**: `/` (This checks the Management UI)
   - **Health Check Port**: `15672`
   - **Disk**: (Recommended) Add a 1GB disk at `/var/lib/rabbitmq` for persistent messages.
8. **Environment Variables**:
   - `RABBITMQ_USER`: Your desired username (e.g., `admin`).
   - `RABBITMQ_PASS`: Your desired password (e.g., `7866`).
   - `RABBITMQ_VHOST`: Your desired vhost (e.g., `barista`).

> [!IMPORTANT]
> **Connecting to RabbitMQ**:
> Other services (like `table` and `qrgen`) must be configured to connect to this instance.
> - **RMQ_HOST**: Use the internal service name with the default AMQP port (e.g., `rabbitmq-service:5672`).
> - **RMQ_USER_NAME**: Must match `RABBITMQ_USER`.
> - **RMQ_PASSWORD**: Must match `RABBITMQ_PASS`.
> - **RMQ_VHOST**: Must match `RABBITMQ_VHOST`.

---

## 5. Environment Variables Reference
Copy `server/.env.example` as a template for your Render environment groups.

### Core Configuration
- `DB_URL`: Ensure `sslmode=require` if using an external database.
- `REDIS_URL`: Use the **Internal Connection String** if Redis is on Render.

---

## 6. Deployment Order
1. **Infrastructure**: Deploy PostgreSQL, Redis, and RabbitMQ (Custom) first.
2. **Internal Services**: Deploy `auth`, `menu`, `shop`, `table`, `history` etc.
3. **Public Services**: Deploy `api-gateway` and `ws` last.

## 7. Health & Readiness Checks
All services are implemented with `cmux` (or standalone HTTP) to expose `/health` and `/ready` on the same port as gRPC. Render uses these to determine if a service should receive traffic or restart.
