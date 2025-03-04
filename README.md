# HextechDronePackets

HextechDronePackets is a Go-based RESTful API for managing drone operations, including tracking drone packets, defining territories, and generating movement reports. It features JWT-based authentication and integrates with a PostgreSQL database, making it suitable for drone telemetry applications.

## Features
- **RESTful API**: Endpoints for user authentication, drone packet submission, territory management, and movement reports.
- **JWT Authentication**: Secures endpoints with JSON Web Tokens.
- **Database Integration**: Uses PostgreSQL for persistent storage.
- **WebSocket Support**: Real-time drone data streaming.
- **Dockerized Deployment**: Easy setup with Docker Compose.

## Getting Started

### Prerequisites
- [Go](https://golang.org/dl/) (1.23+ recommended)
- [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/)
- [cURL](https://curl.se/) (for testing endpoints)
- [wscat](https://www.npmjs.com/package/wscat) (for WebSocket testing, install via `npm install -g wscat`)
- [psql](https://www.postgresql.org/docs/current/app-psql.html) (optional, for database access)

### Clone the Repository
```bash
git clone https://github.com/AVprogpinx131/HextechDronePackets.git
cd HextechDronePackets
```

## Docker setup
Stop existing services (if running):
```bash
docker-compose stop
```
Build and run: 
```bash
docker-compose up -d --build
```

Accessing Postgres database:
```bash
psql -p 5433 -U postgres -h localhost hextech_drone
```


## API endpoints

The API runs on `http://localhost:8080`. Use cURL or similar tools to interact with it. All protected endpoints require a JWT token in the `Authorization: Bearer <token>` header.

### **Authentication**

**Register**: Create a new user.

```bash
curl -X POST http://localhost:8080/register
-H "Content-Type: application/json"
-d '{"username": "example", "password": "password123"}'
```

**Login**: Get a jwt token.
```bash
curl -X POST http://localhost:8080/login
-H "Content-Type: application/json"
-d '{"username": "example", "password": "password123"}'
```


### **Drone packets**

**Submit packet**: Send drone telemetry data.

```bash
curl -X POST http://localhost:8080/api/drone_packet
-H "Content-Type: application/json"
-H "Authorization: <jwt-token>"
-d '{"mac": "AA:BB:CC:DD:EE:FF", "latitude": 40.122, "longitude": 24.237, "altitude": 700}'
```

### **Territories**

**Create**: Define a new territory.
```bash
curl -X POST http://localhost:8080/api/territories 
-H "Content-Type: application/json" 
-H "Authorization: Bearer <jwt-token>" 
-d '{"name": "3d test territory", "latitude": 40.122, "longitude": 24.237, "radius": 500, "min_altitude": 100, "max_altitude": 500}'
```

**List**: Get all user territories.
```bash
curl -X GET http://localhost:8080/api/territories
-H "Content-Type: application/json" 
-H "Authorization: <jwt-token>"
```

**Delete**: Remove a specific territory by ID.
```bash
curl -X DELETE http://localhost:8080/api/territories/{id}
-H "Content-Type: application/json" 
-H "Authorization: Bearer <jwt-token>"
```

### **Reports**

**Drone movements**: View history of tracked drone movements.
```bash
curl -X GET http://localhost:8080/api/reports
-H "Content-Type: application/json" 
-H "Authorization: Bearer <jwt-token>"
```

### **Websocket**

**Connect**: Stream real-time drone data.
```bash
wscat -c "ws://localhost:8080/ws?token=<jwt-token>"
```
NB! Requires the server running locally (`go run cmd/server/main.go`)


## Testing
Integration tests have been made for data and API layer:
```bash
go test ./internal/api -v # API layer
go test ./internal/repository -v # Data layer
```


## Environment Variables

Define these in `config/.env` file:

- `JWT_SECRET`: Secret key for signing JWT tokens.
- `DATABASE_URL`: Connection string for the database.