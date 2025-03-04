# Booking Service

## Overview
The **Booking Service** is a backend application for managing restaurant reservations. It allows users to book tables, and restaurant owners receive notifications via Telegram about new bookings. The service is built using **Golang**, supports **gRPC communication**, and includes a **rate limiter** for API protection.

## Features
- User authentication & authorization (JWT-based)
- Table booking system
- Role-based access control (RBAC)
- Rate limiting to prevent API abuse
- gRPC integration with a **Notification Service** (Telegram bot for restaurant owners)

## Technologies Used
- **Golang** 
- **PostgreSQL** (for storing data)
- **Redis** (for caching)
- **gRPC** (for inter-service communication)
- **JWT** (for authentication)
- **Slog** (for logging)

## Installation & Setup

### 1. Clone the repository
```sh
git clone https://github.com/kourai55k/booking-service
cd booking-service
```

### 2. Build and run the service
```sh
go run cmd/booking-service/main.go
```
