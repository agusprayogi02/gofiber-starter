# Starter Template Go Fiber

Backend Golang With Fiber Framework

### User Guides
1. For Sql Server is not support Enum Data Type Change `entity/user.go` from: 
```go
type User struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Name     string   `gorm:"type:varchar(200);not null"`
	Email    string   `gorm:"type:varchar(200);uniqueIndex;not null"`
	Password string   `gorm:"type:varchar(150);not null"`
	Role     UserRole `gorm:"type:user_role;default:user"`
	gorm.Model
}
```
To :
```go
type User struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Name     string   `gorm:"type:varchar(200);not null"`
	Email    string   `gorm:"type:varchar(200);uniqueIndex;not null"`
	Password string   `gorm:"type:varchar(150);not null"`
	Role     UserRole `gorm:"type:varchar(10);default:user"`
	gorm.Model
}
```
