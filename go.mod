module NYANIMEBACKEND

go 1.23.4

require gorm.io/gorm v1.25.12

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
)

require github.com/go-sql-driver/mysql v1.7.0 // indirect

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/rs/cors v1.11.1
	golang.org/x/crypto v0.31.0
	golang.org/x/text v0.21.0 // indirect
	gorm.io/driver/mysql v1.5.7
)
