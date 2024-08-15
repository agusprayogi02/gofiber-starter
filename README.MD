# Starter Template Go Fiber

Backend Golang With Fiber Framework

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Docker (optional)

### Installation
1. Install `air` for live reload,
```bash
go install github.com/air-verse/air@latest
```
2. Copy `.env.example` to `.env` and fill the values.
3. Change config on .env file, like database connection string.
4. For `Sql Server` change Enum Data Type on `entity/user.go` like: 
```go
type User struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Name     string   `gorm:"type:varchar(200);not null"`
	Email    string   `gorm:"type:varchar(200);uniqueIndex;not null"`
	Password string   `gorm:"type:varchar(150);not null"`
	Role     UserRole `gorm:"type:varchar(10);default:user"` // this
	gorm.Model
}
```
5. Run `air` to start the server with live reload.
6. if you want to run in the docker container, run `docker-compose up -d` and then `docker ps` to check the container is running.
