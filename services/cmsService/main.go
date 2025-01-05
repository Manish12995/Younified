package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"younified-backend/providers/aws"
	"younified-backend/providers/database"
	"younified-backend/providers/graphqlclient"

	controller "younified-backend/services/cmsService/internal/controller"
	resolver "younified-backend/services/cmsService/internal/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const (
	defaultPort         = "4004"
	defaultEnvFile      = ".env"
	defaultDatabaseName = "unified_base"
	defaultRedisHost    = "localhost"
	defaultRedisPort    = 6379
)

// Config holds the application configuration
type Config struct {
	Port          string
	MongoURI      string
	DatabaseName  string
	RedisHost     string
	RedisPort     int
	RedisPassword string
	awsRegion     string
	awsAccessID   string
	awsAccessKey  string
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

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")

	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	awsRegion := os.Getenv("AWS_REGION")

	return Config{
		Port:          port,
		MongoURI:      mongoURI,
		DatabaseName:  defaultDatabaseName,
		RedisHost:     redisHost,
		RedisPort:     redisPort,
		RedisPassword: redisPassword,
		awsRegion:     awsRegion,
		awsAccessID:   awsAccessKeyID,
		awsAccessKey:  awsAccessKey,
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

// initializeGraphQLManager creates a new GraphQL client
func initializeGraphQLManager() *graphqlclient.Graph {
	return graphqlclient.NewGraphql()
}

// initializeRedis sets up the Redis connection

func initializeRedis(config Config) *database.RedisClient {
	redisClient, err := database.NewRedisClient(database.RedisConfig{Host: config.RedisHost, Port: config.RedisPort})
	if err != nil {
		log.Fatalf("Failed to create redis client connection: %v", err)

	}
	return redisClient
}

func initializeAwsService(config Config) *aws.AWSProvider {
	awsProvider, err := aws.NewAWSProvider(config.awsRegion, config.awsAccessID, config.awsAccessKey)
	if err != nil {
		log.Fatalf("Failed to create a session with aws : %v", err)
	}
	return &awsProvider
}

// createGraphQLServer sets up the GraphQL server with resolvers
func createGraphQLServer(
	dbManager *database.DBManager,
	redisClient *database.RedisClient,
	awsProvider *aws.AWSProvider,
	graphqlManager *graphqlclient.Graph) *handler.Server {
	return handler.NewDefaultServer(resolver.NewExecutableSchema(resolver.Config{
		Resolvers: &resolver.Resolver{
			DBManager:     dbManager,
			CMSController: controller.NewCMSController(dbManager, graphqlManager, redisClient, awsProvider),
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
	redisProvider := initializeRedis(config)

	graphqlManager := initializeGraphQLManager()

	awsProvider := initializeAwsService(config)

	// Create GraphQL server
	srv := createGraphQLServer(dbManager, redisProvider, awsProvider, graphqlManager)

	// Setup routes
	setupRoutes(srv)

	// Start server
	startServer(config.Port)
}
