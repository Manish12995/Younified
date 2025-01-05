module younified-backend/services/userService

go 1.23.2

require (
	// Other external dependencies
	github.com/99designs/gqlgen v0.17.56
	go.mongodb.org/mongo-driver v1.17.1
	younified-backend/contracts v0.0.0
	younified-backend/providers/database v0.0.0
	younified-backend/providers/emailBodyProvider v0.0.0
	younified-backend/providers/graphqlclient v0.0.0
)

replace younified-backend/contracts => ../../contracts

replace younified-backend/providers/database => ../../providers/database

replace younified-backend/providers/graphqlclient => ../../providers/graphqlclient

replace younified-backend/providers/emailBodyProvider => ../../providers/emailBodyProvider

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v8 v8.11.5
	github.com/joho/godotenv v1.5.1
	github.com/vektah/gqlparser/v2 v2.5.19
	golang.org/x/crypto v0.27.0
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)
