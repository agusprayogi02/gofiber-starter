# Starter Template Go Fiber

Backend Golang With Fiber Framework

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Docker (optional)
- Open ssl -> untuk generate token jwt

### Installation
1. Install `air` for live reload,
```bash
go install github.com/air-verse/air@latest
```
2. Copy `.env.example` to `.env` and fill the values.
3. Membuat certificate menggunakan open ssl
	- Masuk ke folder `assets/certs` 
	- buat private key
		```bash
		openssl genpkey -algorithm RSA -out certificate.pem -pkeyopt rsa_keygen_bits:4096
		```
	- Masukkan kedalam file `.env` dengan format:
		```bash
		LOCATION_CERT="assets/certs/certificate.pem"
		```
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

## Warning

- Untuk `Order` pada helper digunakan untuk menandai urutan dari error yang muncul.
	- Keterangan kode huruf pertama: 
		- S: Service
		- H: Handler
		- R: Repository
		- U: Util
	- Contoh:
		- S1: Ada pada service dengan error handling pertama

## Aturan
- jika ada error pada `service`, `util`, `helper` wajib mengembalikan error menggunakan `helper` error handling contoh: 
```go
return nil, &helper.BadRequestError{
	Message: err.Error(),
	Order:   "S1",
}
```
- Agar pada `handler` tinggal return error, contoh:
```go
if err != nil {
	return err
}
```

## Rilis Production dengan docker
- pada file `.env` ganti `ENV_TYPE="prod"`
- ganti file `docker-composer.yml` dan `Dockerfile` dengan yang ada didalam folder `assets/docker/prod`, jangan lupa menyesuaikan configurasi yang akan digunakan 


## Authz
- Authz dibuat dengan `casbin` library.
	- https://github.com/gofiber/contrib/blob/main/casbin/casbin_test.go