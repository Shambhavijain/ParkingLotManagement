
#  Parking Lot Management System – CLI & API

This project is a Parking Lot Management System built using **Go**, following **Hexagonal Architecture** and **SOLID principles**. It supports both **CLI-based interaction** and **RESTful APIs** for managing vehicle parking, slot allocation, and ticketing.

---

##  Project Structure

```
parkingSlotManagement/
├── cli/   # CLI entry point (client.go)
|__ cmd/
|   |__app/
|
|__ main.go
|                 
├── internals/
│   ├── adapters/
│   │   └── repositories/   # MySQL & InMemory Repos
│   └── core/
│   |    ├── domain/         # Domain models
│   |    └── services/# Business logic
|   |__ports/   
├── .env                    # Environment variables
└── README.md
```

---

##  Environment Setup

Create a `.env` file in the root directory with the following variables:

```env
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123
JWT_SECRET=your_jwt_secret_key
DB_USER=root
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=3306
DB_NAME=parking_lot
```

---

## Running the CLI

```bash
cd cli
go run client.go
```

Follow the prompts to log in and interact with the system.

---

##  API Endpoints (Testable via Postman)

Base URL: `http://localhost:8080`

| Method | Endpoint              | Description                        |
|--------|-----------------------|------------------------------------|
| POST   | `/login`              | Admin login (returns JWT token)    |
| POST   | `/ParkVehicle`        | Park a vehicle                     |
| POST   | `/UnparkVehicle`      | Unpark a vehicle                   |
| POST   | `/AddSlot`            | Add a new parking slot             |
| GET    | `/GetAvailableSlots`  | View all available slots           |

>  **Note**: Except `/login`, all endpoints require a valid JWT token in the `Authorization` header.

---

##  Sample Postman Request: `/login`

**POST** `http://localhost:8080/login`

**Body (JSON):**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "token": "your-jwt-token"
}
```

Use this token in the `Authorization` header for other requests:

```
Authorization: Bearer your-jwt-token
```

---


##  Technologies Used

- Go (Golang)
- MySQL / InMemory Repositories
- JWT Authentication
- RESTful API
- Postman for testing

---

## Author

Shambhavi Jain

