package controllers

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"younified-backend/contracts/union/model"
	userModel "younified-backend/contracts/user/model"
	"younified-backend/providers/database"
	"younified-backend/providers/graphqlclient"
	"younified-backend/services/unionService/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnionController struct {
	UnionMongoRepository *repository.MongoUnionRepository
	UnionRedisRepository *repository.RedisUnionRepository
	dbManager            *database.DBManager
	graphqlManager       *graphqlclient.Graph
}

func NewUnionController(dbManager *database.DBManager, graphqlManager *graphqlclient.Graph, redisClient *database.RedisClient) *UnionController {
	if dbManager == nil {
		panic("dbManager cannot be nil")
	}

	return &UnionController{
		UnionMongoRepository: repository.NewMongoUnionRepository(dbManager, "unified_base"),
		UnionRedisRepository: repository.NewRedisUnionRepository(redisClient),
		dbManager:            dbManager,
		graphqlManager:       graphqlManager,
	}
}

func (c *UnionController) UnionByID(ctx context.Context, id primitive.ObjectID) (*model.Union, error) {
	// check on cache
	union, err := c.UnionRedisRepository.GetUnionFromCache(ctx, id.Hex())
	if err != nil {
		// get from db, cache and return
		union, err := c.UnionMongoRepository.UnionById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error resolving data from mongodb")
		}
		go c.UnionRedisRepository.CacheUnion(ctx, union.ID.Hex(), union)
		return union, nil
	}
	return union, nil

}

func (c *UnionController) UnionByName(ctx context.Context, name string) (*model.Union, error) {
	space := regexp.MustCompile(`\s+`)
	name = space.ReplaceAllString(strings.TrimSpace(name), " ")
	lower := strings.ToLower(name)
	unionSlug := strings.ReplaceAll(lower, " ", "-")
	nameRegex := `^` + name + `$`
	union, err := c.UnionRedisRepository.GetUnionFromCache(ctx, unionSlug)
	if err != nil {
		// get from db, cache and return
		union, err := c.UnionMongoRepository.UnionByName(ctx, nameRegex)
		if err != nil {
			return nil, fmt.Errorf("error resolving data from mongodb")
		}
		go c.UnionRedisRepository.CacheUnion(ctx, unionSlug, union)
		return union, err
	}
	return union, err
}

func (c *UnionController) Unions(ctx context.Context, page, limit int) (*model.UnionsResponse, error) {

	unions, err := c.UnionRedisRepository.GetUnionsFromCache(ctx, "all-unions") // static key for all unions

	if err != nil {
		// cache miss
		filter := map[string]interface{}{"deleted": false}
		unions, count, err := c.UnionMongoRepository.Unions(ctx, filter, page, limit)
		if err != nil {
			return nil, fmt.Errorf("error resolving data from mongodb")
		}
		// populate cache
		go c.UnionRedisRepository.CacheUnions(ctx, "all-unions", unions)

		return &model.UnionsResponse{
			Unions: unions,
			Count:  int(count),
		}, nil
	}

	return &model.UnionsResponse{
		Unions: unions,
		Count:  int(len(unions)),
	}, nil

}

func (c *UnionController) Register(ctx context.Context, input model.RegisterInput) (*model.Union, error) {
	space := regexp.MustCompile(`\s+`)
	input.Union.UnionID = strings.ToLower(
		space.ReplaceAllString(strings.TrimSpace(input.Union.Name), "-"),
	)
	firstUser := &model.FirstUserInfo{
		FirstName:   input.User.FirstName,
		LastName:    input.User.LastName,
		Email:       input.User.Email,
		Phone:       input.User.Phone,
		Position:    input.User.Position,
		DateOfBirth: input.User.DateOfBirth,
	}

	defaultUser := &model.DefaultUserInfo{
		Username: input.DefaultUser.Username,
		Password: input.DefaultUser.Password,
		Level:    input.DefaultUser.Level,
	}

	union, _ := c.UnionMongoRepository.Register(ctx, &input.Union, firstUser, defaultUser)
	// cache unionByID & unionByName
	go c.UnionRedisRepository.CacheUnion(ctx, union.ID.Hex(), union)
	go c.UnionRedisRepository.CacheUnion(ctx, union.UnionID, union)

	// invalidate cache for all-unions
	_, err := c.UnionRedisRepository.CacheExists(ctx, "all-unions")
	if err == nil {
		// invalidate if cache exists
		go c.UnionRedisRepository.InvalidateCache(ctx, "all-unions")
	}

	mutationInput := map[string]interface{}{
		"unionID":   union.ID,
		"username":  input.Union.FirstUser.Email,
		"password":  input.Union.FirstUser.Password,
		"firstName": input.User.FirstName,
		"lastName":  input.User.LastName,
	}

	gqlEP := os.Getenv("GRAPHQL_ENDPOINT")
	c.graphqlManager.SetgqlEndpoint(gqlEP)
	userMutationBuilder := c.graphqlManager.GetMutationBuilder()

	// Create a user
	createUserMutation, createVars := userMutationBuilder.
		SetMutationName("createUser").
		SetInputName("UserInput").
		SetInput(mutationInput).
		AddField("id").
		Build()
	var result struct {
		RegisterUser *userModel.User `json:"registerUser"`
	}

	go c.graphqlManager.Execute(ctx, createUserMutation, createVars, &result)
	return union, nil
}

func (c *UnionController) ModifyUnion(ctx context.Context, id primitive.ObjectID, union model.Union) (*model.Union, error) {

	updatedUnion, _ := c.UnionMongoRepository.ModifyUnion(ctx, id, union)

	// check cache based on ID
	_, err := c.UnionRedisRepository.CacheExists(ctx, updatedUnion.ID.Hex())
	if err == nil {
		// cache exists so remove it
		go c.UnionRedisRepository.InvalidateCache(ctx, updatedUnion.ID.Hex())
	}
	// check cache based on slug
	_, err = c.UnionRedisRepository.CacheExists(ctx, updatedUnion.UnionID)
	if err == nil {
		// cache exists so remove it
		go c.UnionRedisRepository.InvalidateCache(ctx, updatedUnion.UnionID)
	}

	return updatedUnion, nil
}

func (c *UnionController) DeleteUnion(ctx context.Context, id primitive.ObjectID) (*bool, error) {
	// remove from cache
	go c.UnionRedisRepository.InvalidateCache(ctx, id.Hex())
	// remove from DB
	err := c.UnionMongoRepository.DeleteUnion(ctx, id)
	if err != nil {
		return boolPtr(false), err
	}
	return boolPtr(true), nil
}

func boolPtr(b bool) *bool {
	return &b
}
