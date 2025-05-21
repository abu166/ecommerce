# E-Commerce Microservices Application

## Overview

This project is a microservices-based e-commerce application built using Go, designed to manage products, orders, users, and inventory with a focus on scalability and modularity. The application leverages a microservices architecture, where each service (API Gateway, Inventory, Order, User, Producer, and Consumer) operates independently, communicating via gRPC and NATS for event-driven interactions. It uses PostgreSQL for persistent storage and Docker for containerization, ensuring easy deployment and scalability.

The project demonstrates modern software engineering practices, including Domain-Driven Design (DDD), gRPC for inter-service communication, event-driven architecture with NATS, and container orchestration with Docker Compose.

## Features

The application provides the following functionality:

- **Product Management**: Create, read, update, delete (CRUD) operations for products, including name, category, stock, and price.
- **Order Management**: Create and manage orders, calculate totals based on product prices, and update inventory stock upon order creation.
- **User Management**: User registration, authentication, and profile retrieval with secure password hashing using bcrypt.
- **Inventory Updates**: Real-time inventory updates triggered by order creation using an event-driven approach.
- **API Gateway**: A central entry point for client requests, routing them to appropriate services using RESTful APIs built with Gin.
- **Event-Driven Communication**: Asynchronous communication between services using NATS for publishing and subscribing to order creation events.

## Project Structure

The project follows a modular structure with a focus on Domain-Driven Design (DDD). The directory layout is organized as follows:

```
ecommerce/
├── cmd/                    # Entry points for each microservice
│   ├── apigateway/         # API Gateway service
│   ├── consumer/           # Consumer service for processing order events
│   ├── inventory/          # Inventory management service
│   ├── order/              # Order management service
│   ├── producer/           # Producer service for publishing order events
│   └── user/               # User management service
├── internal/               # Internal packages for each service
│   ├── apigateway/         # API Gateway handlers, middleware, and server logic
│   ├── config/             # Configuration loading from environment variables
│   ├── consumer/           # Consumer service logic
│   ├── inventory/          # Inventory service with DDD layers (application, domain, infrastructure)
│   ├── order/              # Order service with DDD layers
│   ├── producer/           # Producer service logic
│   └── user/               # User service with DDD layers
├── proto/                  # Protocol Buffers (protobuf) definitions for gRPC
├── env                     # Environment configuration file
├── docker-compose.yml      # Docker Compose configuration for orchestration
├── Dockerfile              # Dockerfile for building service images
├── go.mod                  # Go module dependencies
└── README.md               # Project documentation
```

Each service follows a DDD structure with the following layers:

- **Domain**: Defines core business entities (e.g., Product, Order, User) and business logic.
- **Application**: Contains service logic that orchestrates interactions between the domain and infrastructure layers.
- **Infrastructure**: Handles data persistence (PostgreSQL with GORM) and external communication (gRPC, NATS).

## Technologies Used

- **Go**: Primary programming language for all services.
- **gRPC**: For synchronous inter-service communication (e.g., API Gateway to Inventory, Order, User services).
- **NATS**: For asynchronous event-driven communication (e.g., publishing and consuming order creation events).
- **PostgreSQL**: Relational database for storing products, orders, and user data.
- **GORM**: ORM library for database interactions in Go.
- **Gin**: HTTP web framework for the API Gateway's RESTful endpoints.
- **Docker & Docker Compose**: For containerization and orchestration of services.
- **Protocol Buffers (protobuf)**: For defining gRPC service contracts.
- **bcrypt**: For secure password hashing in the User service.
- **godotenv**: For loading environment variables from a .env file.

## Microservices

The application consists of the following microservices:

### API Gateway (cmd/apigateway)

- Exposes RESTful endpoints for clients using Gin.
- Routes requests to Inventory, Order, and User services via gRPC.
- Handles authentication and logging middleware.
- Example endpoints: `/products`, `/orders`, `/users/register`, `/users/login`.

### Inventory Service (cmd/inventory)

- Manages product data (CRUD operations).
- Persists data to PostgreSQL using GORM.
- gRPC service for product-related operations.

### Order Service (cmd/order)

- Manages order creation, retrieval, and updates.
- Checks inventory stock and updates it during order creation.
- Publishes order creation events to NATS via the Producer service.

### User Service (cmd/user)

- Handles user registration, authentication, and profile management.
- Uses bcrypt for password hashing.
- Persists user data to PostgreSQL.

### Producer Service (cmd/producer)

- Publishes order creation events to NATS.
- gRPC service for receiving order notifications from the Order service.

### Consumer Service (cmd/consumer)

- Subscribes to NATS `order.created` events.
- Updates inventory stock based on order items.

## Getting Started with Docker

To run the project using Docker, follow these steps:

### Prerequisites

- Docker and Docker Compose installed on your system.
- Go 1.23 or later (optional, for local development without Docker).
- Basic understanding of microservices and gRPC.

### Steps to Start the Project

1. **Clone the Repository (if applicable)**:
   ```bash
   git clone <repository-url>
   cd ecommerce
   ```

2. **Set Up Environment Variables**:

   Create a `.env` file in the `ecommerce/` directory based on the env file provided:
   ```
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=admin
   DB_NAME=ecommerce
   NATS_ADDR=nats://nats:4222
   API_GATEWAY_ADDR=:8080
   INVENTORY_ADDR=:50051
   ORDER_ADDR=:50052
   USER_ADDR=:50053
   PRODUCER_ADDR=:50054
   ```
   Alternatively, use the provided `env` file.

3. **Build and Run with Docker Compose**:

   Run the following command to start all services (PostgreSQL, NATS, and microservices):
   ```bash
   docker-compose up --build
   ```
   
   This command builds the Docker images for each service and starts the containers.
   The services will be available at:
   - API Gateway: http://localhost:8080
   - Inventory Service: :50051 (gRPC)
   - Order Service: :50052 (gRPC)
   - User Service: :50053 (gRPC)
   - Producer Service: :50054 (gRPC)
   - PostgreSQL: localhost:5432
   - NATS: localhost:4222

4. **Access the Application**:

   Use tools like Postman or cURL to interact with the API Gateway's REST endpoints (e.g., POST /products, GET /orders/:id).
   Example cURL command to create a product:
   ```bash
   curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Laptop","category":"Electronics","stock":100,"price":999.99}'
   ```

5. **Stop the Services**:

   To stop the running containers:
   ```bash
   docker-compose down
   ```

## Database Schema

The application uses PostgreSQL with the following main tables (auto-migrated by GORM):

**Products (inventory service)**:
- `id` (UUID, primary key)
- `name` (string)
- `category` (string)
- `stock` (integer)
- `price` (float64)

**Orders (order service)**:
- `id` (UUID, primary key)
- `user_id` (string)
- `status` (string)
- `total` (float64)
- `created_at`, `updated_at` (timestamps)

**Order Items (order service)**:
- `order_id` (UUID, primary key)
- `product_id` (string)
- `quantity` (integer)

**Users (user service)**:
- `id` (UUID, primary key)
- `username` (string, unique)
- `password` (string, hashed)
- `email` (string, unique)
- `created_at`, `updated_at` (timestamps)

## Event-Driven Workflow

The application uses NATS for event-driven communication:

1. When an order is created, the Order service notifies the Producer service via gRPC.
2. The Producer service publishes an `order.created` event to NATS.
3. The Consumer service subscribes to `order.created` events and updates the inventory stock for each product in the order.

## Dependencies

The project uses the following key Go dependencies (see go.mod for the full list):

- `github.com/gin-gonic/gin`: For RESTful API routing in the API Gateway.
- `google.golang.org/grpc`: For gRPC communication.
- `github.com/nats-io/nats.go`: For NATS messaging.
- `gorm.io/driver/postgres` and `gorm.io/gorm`: For PostgreSQL database interactions.
- `golang.org/x/crypto/bcrypt`: For password hashing.
- `github.com/google/uuid`: For generating UUIDs.

## Future Improvements

- Add unit and integration tests for each service.
- Implement JWT-based authentication for securing API endpoints.
- Add monitoring and logging with tools like Prometheus and Grafana.
- Introduce rate limiting and circuit breaking for resilience.
- Enhance error handling with more detailed error messages.
- Add support for distributed tracing (e.g., Jaeger).

## Notes

- The application assumes a local development environment. For production, update the `.env` file with appropriate values (e.g., database credentials, NATS address).
- Ensure Docker has sufficient resources (CPU, memory) to run all services.
- The API Gateway provides a single entry point, making it easy to extend with additional services.

---

This project showcases my ability to design and implement a scalable microservices architecture using Go, Docker, and modern communication protocols like gRPC and NATS. It demonstrates best practices in software engineering, including DDD, containerization, and event-driven design.
