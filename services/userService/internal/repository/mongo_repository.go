package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"younified-backend/contracts/user/model"
	"younified-backend/providers/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection = "users"
const memberCollection = "members"

type MongoUserRepository struct {
	dbManager  *database.DBManager
	baseDBName string
}

func NewMongoUserRepository(dbManager *database.DBManager, baseDBName string) *MongoUserRepository {
	return &MongoUserRepository{
		dbManager:  dbManager,
		baseDBName: baseDBName,
	}
}

// function to create member while registering on the member collection
func (r *MongoUserRepository) CreateMember(ctx context.Context, unionID string, user *model.User) (*model.User, error) {
	//set time while creating member
	user.CreatedOn = time.Now()
	collection, _ := r.dbManager.GetCollection(ctx, unionID, memberCollection)
	_, err := collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MongoUserRepository) Create(ctx context.Context, unionID string, user *model.User) (*model.User, error) {
	user.CreatedOn = time.Now()

	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	//check this later as the user getting returned is input itself.
	return user, nil
}

func (r *MongoUserRepository) UploadUsers(ctx context.Context, unionID string, bulkWriteModels []mongo.WriteModel, bulkWriteOptions *options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)
	bulkWriteResult, err := collection.BulkWrite(ctx, bulkWriteModels, bulkWriteOptions)
	if err != nil {
		err = fmt.Errorf("could not complete bulk insert %v", err)
		return nil, err
	}
	return bulkWriteResult, nil
}

func (r *MongoUserRepository) GetUser(ctx context.Context, unionID string, filter interface{}) (*model.User, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)
	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) GetByID(ctx context.Context, unionID string, id primitive.ObjectID) (*model.User, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)
	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return	 nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) GetByUsername(ctx context.Context, unionID string, username string) (*model.User, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)
	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": username, "deleted": false}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) GetMemberByID(ctx context.Context, unionID string, id primitive.ObjectID) (*model.User, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, memberCollection)
	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, unionID string, id primitive.ObjectID, updates model.UserUpdateInput) (*model.User, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)

	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedUser model.User
	if err := result.Decode(&updatedUser); err != nil {

		return nil, err
	}
	return &updatedUser, nil
}

func (r *MongoUserRepository) UpdateUser(ctx context.Context, unionID string, filter interface{}, update interface{}) (*model.User, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)

	result := collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var updatedUser model.User
	if err := result.Decode(&updatedUser); err != nil {

		return nil, err
	}
	return &updatedUser, nil
}

func (r *MongoUserRepository) Delete(ctx context.Context, unionID string, id primitive.ObjectID) error {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)

	update := bson.M{
		"$set": bson.M{
			"deleted":   true,
			"deletedAt": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		err = fmt.Errorf("could not delete the user")
		return err
	}

	if result.ModifiedCount == 0 {
		err = fmt.Errorf("user not found")
		return err
	}

	return nil
}

func (r *MongoUserRepository) Restore(ctx context.Context, unionID string, id primitive.ObjectID, update interface{}) error {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		err = fmt.Errorf("could not delete the user")
		return err
	}

	if result.ModifiedCount == 0 {
		err = fmt.Errorf("user not found")
		return err
	}

	return nil
}

func (r *MongoUserRepository) Find(ctx context.Context, unionID string, filter *model.UserFilterInput, page, limit int) ([]*model.User, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)

	findFilter := bson.M{}
	if filter != nil {
		if filter.IsAdmin != false {
			findFilter["isAdmin"] = filter.IsAdmin
		}
		if filter.Deleted != true {
			findFilter["deleted"] = filter.Deleted
		}
		if filter.Status != "" {
			findFilter["status"] = filter.Status
		}
	}

	opts := options.Find()
	if page > 0 && limit > 0 {
		opts.SetSkip(int64((page - 1) * limit))
		opts.SetLimit(int64(limit))
	}

	cursor, err := collection.Find(ctx, findFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *MongoUserRepository) Count(ctx context.Context, unionID string, filter *model.UserFilterInput) (int64, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, userCollection)

	findFilter := bson.M{}
	if filter != nil {
		if filter.IsAdmin != false {
			findFilter["isAdmin"] = filter.IsAdmin
		}
		if filter.Deleted != true {
			findFilter["deleted"] = filter.Deleted
		}
		if filter.Status != "" {
			findFilter["status"] = filter.Status
		}
	}

	return collection.CountDocuments(ctx, findFilter)
}
