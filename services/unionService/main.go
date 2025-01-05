package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"younified-backend/providers/database"
	"younified-backend/providers/graphqlclient"
	controllers "younified-backend/services/unionService/internal/controller"
	resolver "younified-backend/services/unionService/internal/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const (
	defaultPort         = "4001"
	defaultRedisPort    = 6379
	defaultRedisHost    = "localhost"
	defaultEnvFile      = ".env"
	defaultDatabaseName = "unified_base"
)

// Config holds the application configuration
type Config struct {
	Port          string
	MongoURI      string
	RedisHost     string
	RedisPort     int
	RedisPassword string
	DatabaseName  string
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

	redisPortStr := os.Getenv("REDIS_PORT")
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		redisPort = defaultRedisPort
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	return Config{
		Port:          port,
		MongoURI:      mongoURI,
		RedisHost:     redisHost,
		RedisPort:     redisPort,
		RedisPassword: redisPassword,
		DatabaseName:  defaultDatabaseName,
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

func initializeRedis(config Config) *database.RedisClient {
	redisClient, err := database.NewRedisClient(database.RedisConfig{Host: config.RedisHost, Port: config.RedisPort})
	if err != nil {
		log.Fatalf("Failed to create redis client connection: %v", err)
	}
	return redisClient
}

// initializeGraphQLManager creates a new GraphQL client
func initializeGraphQLManager() *graphqlclient.Graph {
	return graphqlclient.NewGraphql()
}

// createGraphQLServer sets up the GraphQL server with resolvers
func createGraphQLServer(
	dbManager *database.DBManager,
	graphqlManager *graphqlclient.Graph,
	redisClient *database.RedisClient,
) *handler.Server {
	return handler.NewDefaultServer(resolver.NewExecutableSchema(resolver.Config{
		Resolvers: &resolver.Resolver{
			DBManager:       dbManager,
			UnionController: controllers.NewUnionController(dbManager, graphqlManager, redisClient),
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

	// Initialize dependencies
	dbManager := initializeDatabase(ctx, config)
	defer dbManager.Close(ctx)

	redisClient := initializeRedis(config)
	defer redisClient.Close()

	graphqlManager := initializeGraphQLManager()

	// Create GraphQL server
	srv := createGraphQLServer(dbManager, graphqlManager, redisClient)

	// Setup routes
	setupRoutes(srv)

	// Start server
	startServer(config.Port)
}
