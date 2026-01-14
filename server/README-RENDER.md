# Deploying Barista Microservices to Render

This guide provides step-by-step instructions for deploying the Barista project as individual services on [Render](https://render.com). The architecture separates **Public Web Services** (accessible from the internet) from **Private Services** (internal gRPC/AMQP services).

## 1. Prerequisites
- A Render account.
- A GitHub repository with the Barista code.
- A PostgreSQL database (Render managed or external).
- A Redis instance (Render managed).
- A RabbitMQ instance (Render managed or external).

---

## 2. Public Web Services
These services will have a public URL (e.g., `https://api-gateway-xxx.onrender.com`).

### api-gateway
1. Click **New +** > **Web Service**.
2. Connect your repository.
3. **Region**: Choose the same region for all services to minimize latency.
4. **Build Command**: `go build -o server ./api-gateway/cmd/main.go`
5. **Start Command**: `./server`
6. **Environment Variables**:
   - `PORT`: `8080` (Render will map traffic to this)
   - `JWT_KEY`: A long random string.
   - `ALLOWED_ORIGINS`: Your frontend URL.
   - `AUTH_SERVICE_URL`, `SHOP_SERVICE_URL`, etc.: Use the **Internal Service Name** (e.g., `http://auth-service:8080`).

### ws (WebSocket Service)
1. Click **New +** > **Web Service**.
2. **Build Command**: `go build -o server ./services/ws/cmd/main.go`
3. **Environment Variables**:
   - `PORT`: `8080`
   - `ALLOWED_ORIGINS`: Your frontend URL.

---

## 3. Private Services
These services are **not** accessible from the internet. They are reachable only via internal hostnames provided by Render.

### General Steps for Internal Services
(`auth`, `order`, `payment`, `session`, `shop`, `table`, `history`, `menu`, `qrgen`)

1. Click **New +** > **Private Service**.
2. Connect your repository.
3. **Build Command**: `go build -o server ./services/<service_name>/cmd/main.go`
4. **Start Command**: `./server`
5. **Advanced Settings**:
   - **Health Check Path**: `/health`
6. **Environment Variables**:
   - `PORT`: `8080` (Internal port)
   - `DB_URL`: PostgreSQL connection string.
   - `REDIS_URL`: (For `session` service) External/Internal Redis string.
   - `*_SERVICE_URL`: Internal URLs for other dependencies (e.g., `http://menu-service:8080`).

> [!TIP]
> **Internal Hostnames**: When service `A` needs to call service `B`, use `http://service-b-name:8080`. You can find the internal name in the service's dashboard.

---

## 4. Environment Variables Reference
Copy `server/.env.example` as a template for your Render environment groups.

### Core Configuration
- `DB_URL`: Ensure `sslmode=require` if using an external database.
- `REDIS_URL`: Use the **Internal Connection String** if Redis is on Render.

---

## 5. Deployment Order
1. **Infrastructure**: Deploy PostgreSQL, Redis, and RabbitMQ first.
2. **Internal Services**: Deploy `auth`, `menu`, `shop`, `table`, `history` etc.
3. **Public Services**: Deploy `api-gateway` and `ws` last.

## 6. Health & Readiness Checks
All services are implemented with `cmux` (or standalone HTTP) to expose `/health` and `/ready` on the same port as gRPC. Render uses these to determine if a service should receive traffic or restart.
