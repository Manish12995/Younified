package repository

import (
	"context"
	"errors"
	"time"
	union "younified-backend/contracts/union/model"
	"younified-backend/providers/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUnionRepository struct {
	dbManager  *database.DBManager
	baseDBName string
}

func NewMongoUnionRepository(dbManager *database.DBManager, baseDBName string) *MongoUnionRepository {
	return &MongoUnionRepository{
		dbManager:  dbManager,
		baseDBName: baseDBName,
	}
}

func (r *MongoUnionRepository) Register(ctx context.Context, unionData *union.Union, firstUser *union.FirstUserInfo, defaultUser *union.DefaultUserInfo) (*union.Union, error) {
	// Use unified_base database
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")

	// Set creation timestamp
	unionData.CreatedOn = time.Now()
	unionData.FirstUser = *firstUser
	unionData.DefaultUser = *defaultUser

	// Insert the union
	result, err := unionCollection.InsertOne(ctx, unionData)
	if err != nil {
		return nil, err
	}

	// Set the generated ID
	unionData.ID = result.InsertedID.(primitive.ObjectID)
	return unionData, nil
}

// DeleteUnion permanently removes a Union from the database
func (r *MongoUnionRepository) DeleteUnion(ctx context.Context, unionID primitive.ObjectID) error {
	// Perform a hard delete
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	filter := bson.M{"_id": unionID}
	result, err := unionCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no union found with the given ID")
	}

	return nil
}

// ModifyUnion updates an existing Union's information
func (r *MongoUnionRepository) ModifyUnion(ctx context.Context, unionID primitive.ObjectID, update union.Union) (*union.Union, error) {
	// Prepare the update document
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	updateDoc := bson.M{
		"$set": update,
	}

	// Update options to return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// Perform the update
	var updatedUnion union.Union
	err := unionCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": unionID},
		updateDoc,
		opts,
	).Decode(&updatedUnion)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no union found with the given ID")
		}
		return nil, err
	}

	return &updatedUnion, nil
}

// UnionById retrieves a Union by its unique identifier
func (r *MongoUnionRepository) UnionById(ctx context.Context, unionID primitive.ObjectID) (*union.Union, error) {
	// Use read preference if set
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	var findOptions *options.FindOneOptions

	// Find the union
	var foundUnion union.Union
	err := unionCollection.FindOne(ctx, bson.M{"_id": unionID}, findOptions).Decode(&foundUnion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no union found with the given ID")
		}
		return nil, err
	}

	return &foundUnion, nil
}

// UnionByName retrieves a Union by its name
func (r *MongoUnionRepository) UnionByName(ctx context.Context, name string) (*union.Union, error) {
	// Use read preference if set
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	var findOptions *options.FindOneOptions

	// Find the union
	var foundUnion union.Union
	err := unionCollection.FindOne(ctx, bson.M{"name": primitive.Regex{Pattern: name, Options: "i"}, "deleted": false}, findOptions).Decode(&foundUnion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no union found with the given name")
		}
		return nil, err
	}

	return &foundUnion, nil
}

// Unions retrieves all Unions with optional filtering and pagination
func (r *MongoUnionRepository) Unions(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*union.Union, int64, error) {
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	// Prepare find options with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * limit))
	findOptions.SetLimit(int64(limit))

	// Convert filter to BSON
	mongoBSONFilter := bson.M(filter)

	// Count total documents for pagination
	totalCount, err := unionCollection.CountDocuments(ctx, mongoBSONFilter)
	if err != nil {
		return nil, 0, err
	}

	// Find the unions
	cursor, err := unionCollection.Find(ctx, mongoBSONFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var unions []*union.Union
	if err = cursor.All(ctx, &unions); err != nil {
		return nil, 0, err
	}

	return unions, totalCount, nil
}

// Exists checks if a Union exists by a given condition
func (r *MongoUnionRepository) Exists(ctx context.Context, filter map[string]interface{}) (bool, error) {
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	// Convert filter to BSON
	mongoBSONFilter := bson.M(filter)

	// Use read preference if set
	var findOptions *options.CountOptions

	// Count documents matching the filter
	count, err := unionCollection.CountDocuments(ctx, mongoBSONFilter, findOptions)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Additional methods from the previous interface...
func (r *MongoUnionRepository) AddModule(ctx context.Context, unionID primitive.ObjectID, module string) error {
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")

	filter := bson.M{"_id": unionID}
	update := bson.M{
		"$addToSet": bson.M{"modules": module},
	}

	_, err := unionCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoUnionRepository) RemoveModule(ctx context.Context, unionID primitive.ObjectID, module string) error {
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	filter := bson.M{"_id": unionID}
	update := bson.M{
		"$pull": bson.M{"modules": module},
	}

	_, err := unionCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoUnionRepository) UpdateDefaultUser(ctx context.Context, unionID primitive.ObjectID, defaultUser *union.DefaultUserInfo) error {
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	filter := bson.M{"_id": unionID}
	update := bson.M{
		"$set": bson.M{"defaultUser": defaultUser},
	}

	_, err := unionCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoUnionRepository) SetManagerIDs(ctx context.Context, unionID primitive.ObjectID, accountManagers, communicationReps []primitive.ObjectID) error {
	unifiedDB := r.dbManager.GetBaseDatabase(ctx)
	unionCollection := unifiedDB.Collection("unions")
	filter := bson.M{"_id": unionID}
	update := bson.M{
		"$set": bson.M{
			"accountManager":   accountManagers,
			"communicationRep": communicationReps,
		},
	}

	_, err := unionCollection.UpdateOne(ctx, filter, update)
	return err
}
