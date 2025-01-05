package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client *mongo.Client
	once   sync.Once
)

// Config holds MongoDB connection configuration
type Config struct {
	URI      string
	Database string
	Username string
	Password string
	// Add other configuration options as needed
}

// MongoDB represents the MongoDB client wrapper
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB creates a new MongoDB instance with the given configuration
func NewMongoDB(cfg *Config) (*MongoDB, error) {
	var err error

	// Use sync.Once to ensure single instance
	once.Do(func() {
		// Create client options
		clientOptions := options.Client().
			ApplyURI(cfg.URI).
			SetMaxPoolSize(100).
			SetMinPoolSize(10).
			SetMaxConnIdleTime(30 * time.Second).
			SetReadPreference(readpref.SecondaryPreferred()).
			SetRetryWrites(true)

		// Add credentials if provided
		if cfg.Username != "" && cfg.Password != "" {
			credential := options.Credential{
				Username: cfg.Username,
				Password: cfg.Password,
			}
			clientOptions.SetAuth(credential)
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			return
		}

		// Verify connection
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Printf("Failed to ping MongoDB: %v", err)
			return
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB: %v", err)
	}

	// Create MongoDB instance
	mongodb := &MongoDB{
		Client:   client,
		Database: client.Database(cfg.Database),
	}

	return mongodb, nil
}

// Collection returns a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

// Disconnect closes the MongoDB connection
func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// HealthCheck verifies if the MongoDB connection is healthy
func (m *MongoDB) HealthCheck(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

// WithTransaction executes the provided function within a transaction
func (m *MongoDB) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := m.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	wrappedFn := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := fn(sessCtx)
		return nil, err // Return nil as interface{} and the error
	}

	_, err = session.WithTransaction(ctx, wrappedFn)
	return err
}
