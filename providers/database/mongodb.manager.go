package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Union struct {
	ID      primitive.ObjectID `bson:"_id"`
	UnionID string             `bson:"unionID"`
}

type DBManager struct {
	client         *mongo.Client
	databases      map[string]*mongo.Database
	dbNameMap      map[string]string // Maps ObjectID to actual database name
	mu             sync.RWMutex
	baseDBName     string
	serviceDBNames map[string]string
}

func NewDBManager(ctx context.Context, uri string, baseDBName string) (*DBManager, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30*time.Second).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetRetryWrites(true))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	return &DBManager{
		client:         client,
		databases:      make(map[string]*mongo.Database),
		dbNameMap:      make(map[string]string),
		baseDBName:     baseDBName,
		serviceDBNames: make(map[string]string),
	}, nil
}

func (m *DBManager) FindDBName(ctx context.Context, key string) (string, error) {
	// First check the cache
	m.mu.RLock()
	if dbName, exists := m.dbNameMap[key]; exists {
		m.mu.RUnlock()
		return dbName, nil
	}
	m.mu.RUnlock()

	// Check dynamically set database names first
	if dbName, exists := m.GetServiceDBName(key); exists {
		return dbName, nil
	}

	// If not a dynamic name, try union database lookup
	objID, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return "", fmt.Errorf("invalid unionID: %v", err)
	}
	unifiedDB := m.client.Database("unified_base")
	var union Union
	err = unifiedDB.Collection("unions").FindOne(ctx, bson.M{"_id": objID}).Decode(&union)
	if err != nil {
		return "", fmt.Errorf("failed to find union: %v", err)
	}
	m.mu.Lock()
	m.dbNameMap[key] = union.UnionID
	m.mu.Unlock()

	return union.UnionID, nil
}

func (m *DBManager) GetDatabase(ctx context.Context, key string) (*mongo.Database, error) {
	m.mu.RLock()
	if db, exists := m.databases[key]; exists {
		m.mu.RUnlock()
		return db, nil
	}

	m.mu.RUnlock()

	dbName, err := m.FindDBName(ctx, key)
	if err != nil {
		return nil, err
	}
	fmt.Println("print db: ", dbName)
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if db, exists := m.databases[key]; exists {
		return db, nil
	}

	db := m.client.Database(dbName)
	m.databases[key] = db
	return db, nil
}
func (m *DBManager) GetCollection(ctx context.Context, unionID, collectionName string) (*mongo.Collection, error) {
	db, err := m.GetDatabase(ctx, unionID)
	if err != nil {
		return nil, err
	}
	return db.Collection(collectionName), nil
}

func (m *DBManager) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *DBManager) GetBaseDatabase(ctx context.Context) *mongo.Database {
	return m.client.Database(m.baseDBName)
}

// SetDynamicDBName allows dynamic setting of database names
func (m *DBManager) SetServiceDBName(key string, dbName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.serviceDBNames[key] = dbName
}

// GetDynamicDBName retrieves a dynamically set database name
func (m *DBManager) GetServiceDBName(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	dbName, exists := m.serviceDBNames[key]
	return dbName, exists
}
