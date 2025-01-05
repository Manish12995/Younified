package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"younified-backend/contracts/cms/model"
	"younified-backend/providers/database"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const BlogServiceDBKey = "blogs"
const BlogServiceDBName = "blog"
const BlogServiceCollection = "blog-post"

type MongoCommsRepository struct {
	dbManager  *database.DBManager
	baseDBName string
}

func NewMongoCommsRepository(dbManager *database.DBManager, baseDBName string) *MongoCommsRepository {
	fmt.Println("check the setService call")
	dbManager.SetServiceDBName(BlogServiceDBKey, BlogServiceDBName)
	return &MongoCommsRepository{
		dbManager:  dbManager,
		baseDBName: baseDBName,
	}

}

func (r *MongoCommsRepository) CreateNews(ctx context.Context, unionID string, input model.News) (*model.News, error) {
	newsCollection, _ := r.dbManager.GetCollection(ctx, unionID, "news")
	insertedID, err := newsCollection.InsertOne(ctx, input)
	if err != nil {
		return nil, err
	}
	news := new(model.News)
	newsResult := newsCollection.FindOne(ctx, bson.M{"_id": &insertedID.InsertedID})
	err = newsResult.Decode(news)
	if err != nil {
		logrus.Infof("could not decode inserted news %v", news)
		err = fmt.Errorf("could not decode inserted news : %v", err)
		return nil, err
	}
	return news, nil
}

// Unions retrieves all Unions with optional filtering and pagination
func (r *MongoCommsRepository) GetAllNewsPosts(ctx context.Context, unionID string, filter interface{}, page, limit int, sort interface{}) ([]*model.News, int64, error) {
	newsCollection, _ := r.dbManager.GetCollection(ctx, unionID, "news")
	// Prepare find options with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * limit))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(sort)

	// Convert filter to BSON

	// Count total documents for pagination
	totalCount, err := newsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	skip := (page - 1) * limit
	// Find the unions
	pipeline := []bson.M{{
		"$lookup": bson.M{
			"from":         "users",
			"localField":   "userID",
			"foreignField": "_id",
			"as":           "creator",
		}},
		{
			"$unwind": bson.M{
				"path":                       "$creator",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "likes",
				"foreignField": "_id",
				"as":           "likedBy",
			},
		},
		{"$match": filter},
		{"$sort": sort},
		{"$skip": skip},
		{"$limit": limit},
	}

	opts := options.Aggregate()
	allowDiskUse := true
	opts.AllowDiskUse = &allowDiskUse

	cursor, err := newsCollection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		err = fmt.Errorf("could not fetch news data")
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var news []*model.News
	if err = cursor.All(ctx, &news); err != nil {
		return nil, 0, err
	}

	return news, totalCount, nil
}

func (r *MongoCommsRepository) GetBlogs(ctx context.Context, filter interface{}) ([]*model.Blog, error) {
	collection, err := r.dbManager.GetCollection(ctx, BlogServiceDBKey, BlogServiceCollection)
	if err != nil {
		log.Printf("could not fetch collection : %v", err)
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		err = fmt.Errorf("could not fetch news data")
		return nil, err
	}
	defer cursor.Close(ctx)
	var blogs []*model.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *MongoCommsRepository) GetOneBlog(ctx context.Context, filter interface{}) (*model.Blog, error) {
	collection, err := r.dbManager.GetCollection(ctx, BlogServiceDBKey, BlogServiceCollection)
	if err != nil {
		log.Printf("could not fetch collection : %v", err)
	}
	result := collection.FindOne(ctx, filter)
	blog := new(model.Blog)
	err = result.Decode(&blog)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (r *MongoCommsRepository) GetOneNews(ctx context.Context, unionID string, newsID primitive.ObjectID) (*model.News, error) {
	collection, _ := r.dbManager.GetCollection(ctx, unionID, "news")
	var findOptions *options.FindOneOptions

	// Find the union
	var foundNews model.News
	err := collection.FindOne(ctx, bson.M{"_id": newsID}, findOptions).Decode(&foundNews)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no union found with the given ID")
		}
		return nil, err
	}

	return &foundNews, nil
}

func (r *MongoCommsRepository) GetOneComment(ctx context.Context, unionID string, newsID primitive.ObjectID, commentID primitive.ObjectID) (*model.Comment, error) {
	commentCollection, _ := r.dbManager.GetCollection(ctx, unionID, "newscomments_"+newsID.Hex())
	var findOptions *options.FindOneOptions

	// Find the union
	var foundComment model.Comment
	err := commentCollection.FindOne(ctx, bson.M{"_id": commentID}, findOptions).Decode(&foundComment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no union found with the given ID")
		}
		return nil, err
	}

	return &foundComment, nil
}

func (r *MongoCommsRepository) RetrieveComments(ctx context.Context, filter interface{}, unionID string, newsID string, page int, limit int) ([]*model.Comment, error) {

	commentCollection, _ := r.dbManager.GetCollection(ctx, unionID, "newscomments_"+newsID)
	skip := (page - 1) * limit
	sort := bson.M{"createdOn": -1}
	pipeline := []bson.M{{
		"$lookup": bson.M{
			"from":         "users",
			"localField":   "userID",
			"foreignField": "_id",
			"as":           "creator",
		}},
		{
			"$unwind": bson.M{
				"path":                       "$creator",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{"$sort": bson.M{
			"createdOn": -1,
		}},
		{"$match": filter},
		{"$skip": skip},
		{"$sort": sort},
	}
	opts := options.Aggregate()
	allowDiskUse := true
	opts.AllowDiskUse = &allowDiskUse

	// logrus.Tracef("created aggregate pipeline %+v", pipeline)

	cursor, err := commentCollection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		err = fmt.Errorf("could not fetch comments data")
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var comments []*model.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *MongoCommsRepository) UpdateNews(ctx context.Context, unionID string, filter interface{}, update interface{}) (*model.News, error) {
	fmt.Println("called updateNews")
	collection, _ := r.dbManager.GetCollection(ctx, unionID, "news")

	result := collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	fmt.Println("", *result)
	var updatedNews model.News
	if err := result.Decode(&updatedNews); err != nil {

		return nil, err
	}
	return &updatedNews, nil
}

func (r *MongoCommsRepository) AddComment(ctx context.Context, unionID string, newsID primitive.ObjectID, input *model.Comment) (*model.Comment, error) {
	commentCollection, _ := r.dbManager.GetCollection(ctx, unionID, "newscomments_"+newsID.Hex())
	logrus.Tracef("Initializing insertion of comment")
	insertedID, err := commentCollection.InsertOne(ctx, input)
	if err != nil {
		return nil, err
	}
	var comment model.Comment
	err = commentCollection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(comment)
	if err != nil {
		logrus.Tracef("could not decode inserted comment %v", comment)
		err = fmt.Errorf("could not decode inserted comment")
		return nil, err
	}
	return &comment, nil
}

func (r *MongoCommsRepository) UpdateOneComment(ctx context.Context, unionID string, newsID primitive.ObjectID, filter interface{}, update interface{}) (*model.Comment, error) {
	commentCollection, _ := r.dbManager.GetCollection(ctx, unionID, "newscomments_"+newsID.Hex())
	logrus.Tracef("Initializing update of a comment")
	result := commentCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var UpdatedComment model.Comment
	if err := result.Decode(&UpdatedComment); err != nil {
		return nil, err
	}
	return &UpdatedComment, nil
}

func (r *MongoCommsRepository) CreateBlog(ctx context.Context, input model.Blog) (*model.Blog, error) {
	blogsCollection, _ := r.dbManager.GetCollection(ctx, BlogServiceDBKey, BlogServiceCollection)
	_, err := blogsCollection.InsertOne(ctx, input)

	if err != nil {
		err = fmt.Errorf("error: %V", err)
		return nil, err
	}
	return &input, nil
}

func (r *MongoCommsRepository) UpdateBlog(ctx context.Context, filter interface{}, input interface{}) (*model.Blog, error) {
	blogsCollection, _ := r.dbManager.GetCollection(ctx, BlogServiceDBKey, BlogServiceCollection)

	result := blogsCollection.FindOneAndUpdate(
		ctx,
		filter,
		input,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var updatedBlog model.Blog
	if err := result.Decode(&updatedBlog); err != nil {

		return nil, err
	}
	return &updatedBlog, nil
}

func (r *MongoCommsRepository) DeleteBlog(ctx context.Context, filter interface{}) (*model.Blog, error) {
	blogsCollection, _ := r.dbManager.GetCollection(ctx, BlogServiceDBKey, BlogServiceCollection)
	updateDoc := bson.M{
		"$set": bson.M{"deleted": true},
	}
	result := blogsCollection.FindOneAndUpdate(
		ctx,
		filter,
		updateDoc,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var updatedBlog model.Blog
	if err := result.Decode(&updatedBlog); err != nil {

		return nil, err
	}
	return &updatedBlog, nil
}
