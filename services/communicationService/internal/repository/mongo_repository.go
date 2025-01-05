package repository

import "younified-backend/providers/database"

type MongoCommsRepository struct {
	dbManager  *database.DBManager
	baseDBName string
}

func NewMongoCommsRepository(dbManager *database.DBManager, baseDBName string) *MongoCommsRepository {
	return &MongoCommsRepository{
		dbManager:  dbManager,
		baseDBName: baseDBName,
	}
}
