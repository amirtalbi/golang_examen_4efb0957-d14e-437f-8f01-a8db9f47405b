Examen
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   ├── health_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/
│   │   │   └── auth_middleware.go
│   │   └── routes/
│   │       └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── models/
│   │   │   └── user.go
│   │   └── repositories/
│   │       └── user_repository.go
│   └── service/
│       ├── auth_service.go
│       └── user_service.go
├── pkg/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   └── validator/
│       └── validator.go
├── .env
├── .env.example
├── .github/
│   └── workflows/
│       └── ci.yml
├── go.mod
├── go.sum
└── README.md