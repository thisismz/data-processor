# data processor in golang

## Description

This project implements a Data Processor system with an input component, a processing queue, and a storage space at the end of the queue.
The system follows a Distributed Domain-Driven Design (DDD) architecture and utilizes design patterns such as Circuit Breaker.

## Key Features
- Accepts data with unique identifiers and prevents processing duplicate data
- Enforces rate limiting and traffic quotas for each user, including requests per minute and total data volume per month
- Allows for multiple instances of the service to handle high request volumes
- Configurable to run the service as either a Consumer or a Producer
- Utilizes Docker for easy deployment and scalability


## Tech Stack

- Language: Go
- Storage: Redis, MySQL
- Message Queue: RabbitMQ
- Architecture: Distributed Domain-Driven Design (DDD)
- Design Pattern: Circuit Breaker
- Containerization: Docker

## Getting Started

1. Clone the repository
2. Configure the environment variables (see example.env file)
3. Build and run the application using Docker Compose: docker-compose up -d
### Run web server

```bash
go run cmd/web/main.go
```

## API Spec

All API Spec is in `api` folder. and swagger goto : http://127.0.0.1:80/docs/swagger/

## Database Migration

All database migration is in `migration` folder.

### Create Migration

```shell
migrate create -ext sql -dir db/migrations create_table_xxx
```

### Run Migration

```shell
migrate -database "mysql://root:@tcp(localhost:3306)/?charset=utf8mb4&parseTime=True&loc=Local" -path migration up
```
