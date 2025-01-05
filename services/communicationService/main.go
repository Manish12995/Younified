package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"younified-backend/providers/database"
	controller "younified-backend/services/communicationService/internal/controller"
	resolver "younified-backend/services/communicationService/internal/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const (
	defaultPort         = "4003"
	defaultEnvFile      = ".env"
	defaultDatabaseName = "unified_base"
	defaultRedisHost    = "localhost"
	defaultRedisPort    = 6379
)

// Config holds the application configuration
type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

// loadConfiguration reads environment variables and returns a Config
func loadConfiguration() Config {
	// Load .env file
	if err := godotenv.Load(defaultEnvFile); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI must be set")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = defaultRedisHost
	}

	// redisPortStr := os.Getenv("REDIS_PORT")
	// redisPort, err := strconv.Atoi(redisPortStr)
	// if err != nil {
	// 	redisPort = defaultRedisPort
	// }

	// redisPassword := os.Getenv("REDIS_PASSWORD")

	return Config{
		Port:         port,
		MongoURI:     mongoURI,
		DatabaseName: defaultDatabaseName,
		// RedisHost:     redisHost,
		// RedisPort:     redisPort,
		// RedisPassword: redisPassword,
	}
}

// initializeDatabase sets up the MongoDB connection
func initializeDatabase(ctx context.Context, config Config) *database.DBManager {
	dbManager, err := database.NewDBManager(ctx, config.MongoURI, config.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to create DB manager: %v", err)
	}
	return dbManager
}

// // initializeRedis sets up the Redis connection
// func initializeRedis(config Config) *redis.Provider {
// 	redisConfig := redis.Config{
// 		Host:     config.RedisHost,
// 		Port:     config.RedisPort,
// 		Password: config.RedisPassword,
// 	}
// 	return redis.NewProvider(redisConfig)
// }

// createGraphQLServer sets up the GraphQL server with resolvers
func createGraphQLServer(
	dbManager *database.DBManager,
	// redisProvider *redis.Provider,
) *handler.Server {
	return handler.NewDefaultServer(resolver.NewExecutableSchema(resolver.Config{
		Resolvers: &resolver.Resolver{
			DBManager:       dbManager,
			CommsController: controller.NewCommsController(dbManager),
		},
	}))
}

// setupRoutes configures HTTP routes
func setupRoutes(srv *handler.Server) {
	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)
}

// startServer begins listening on the specified port
func startServer(port string) {
	log.Printf("Connecting to GraphQL playground at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	// Create a context
	ctx := context.Background()

	// Load configuration
	config := loadConfiguration()

	// Initialize database
	dbManager := initializeDatabase(ctx, config)
	defer dbManager.Close(ctx)

	// Initialize Redis
	// redisProvider := initializeRedis(config)

	// Create GraphQL server
	srv := createGraphQLServer(dbManager)

	// Setup routes
	setupRoutes(srv)

	// Start server
	startServer(config.Port)
}
