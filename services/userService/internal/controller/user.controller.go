package controllers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"younified-backend/contracts/user/model"
	"younified-backend/providers/database"
	email "younified-backend/providers/emailBodyProvider"
	"younified-backend/providers/graphqlclient"
	"younified-backend/services/userService/internal/auth"
	"younified-backend/services/userService/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserController struct {
	UserMongoRepository *repository.MongoUserRepository
	UserRedisRepository *repository.RedisUserRepository
	dbManager           *database.DBManager
	graphqlManager      *graphqlclient.Graph
}

func NewUserController(dbManager *database.DBManager, graphqlManager *graphqlclient.Graph, redisClient *database.RedisClient) *UserController {
	if dbManager == nil {
		panic("dbManager cannot be nil")
	}
	return &UserController{
		UserMongoRepository: repository.NewMongoUserRepository(dbManager, "unified_base"),
		UserRedisRepository: repository.NewRedisUserRepository(redisClient),
		dbManager:           dbManager,
		graphqlManager:      graphqlManager,
	}
}

var Response string = "Operation Successful"

//	----------------------- --------------------- -------------------- MAKERS --------------------- ------------------------------- --------------------------
//
// function to create user from interservice communication
func (c *UserController) CreateUser(ctx context.Context, input model.User) (*model.User, error) {
	if !auth.IsPasswordCompromised(input.Password) {
		err := errors.New("password doesn't match the required criteria")
		return nil, err
	}
	// hash the password
	password, _ := auth.HashPassword(input.Password, input.UnionID.Hex())
	var level int = 1
	var deleted bool = false
	var isAdmin bool = false

	if input.Level == 5 {
		// this is default user
		level = 5
		deleted = true
		isAdmin = true
	}
	user := &model.User{
		UnionID:   input.UnionID,
		Username:  input.Username,
		Password:  password, // Add user password hashing strategy here
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Profile:   input.Profile,
		Deleted:   deleted,
		Level:     level,
		Status:    "active",
		IsAdmin:   isAdmin,
	}

	user, _ = c.UserMongoRepository.Create(ctx, input.UnionID.Hex(), user)
	// cache it
	go c.UserRedisRepository.CacheUser(ctx, user.ID.Hex(), user)

	return user, nil
}

// function to register user(member) from client side
func (c *UserController) CreateMember(ctx context.Context, input model.User) (*model.User, error) {
	if !auth.IsPasswordCompromised(input.Password) {
		err := errors.New("password doesn't match the required criteria")
		return nil, err
	}
	unionID := input.UnionID.Hex()
	// hash the password
	password, _ := auth.HashPassword(input.Password, unionID)
	var level int = 1
	var deleted bool = false
	var isAdmin bool = false

	if input.Level == 5 {
		// this is default user
		level = 5
		deleted = true
		isAdmin = true
	}

	user := &model.User{
		UnionID:   input.UnionID,
		Username:  input.Username,
		Password:  password, // Add user password hashing strategy here
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Profile:   input.Profile,
		Deleted:   deleted,
		Level:     level,
		Status:    "registered",
		IsAdmin:   isAdmin,
	}

	// no need to cache memeber registration requests - on approval cache it
	return c.UserMongoRepository.CreateMember(ctx, unionID, user)
}

func (c *UserController) UploadUsers(ctx context.Context, unionID primitive.ObjectID, input []*model.User) (*string, error) {
	if unionID.IsZero() {
		err := fmt.Errorf("unionID is required")
		return nil, err
	}
	var bulkWriteModels []mongo.WriteModel
	for _, user := range input {
		// Prepare the filter and update
		filter := bson.M{"username": user.Username}
		update := bson.M{
			"$set": user,
		}
		// Create an upsert operation
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		bulkWriteModels = append(bulkWriteModels, model)
	}
	// Perform bulk write
	bulkWriteOptions := options.BulkWrite().SetOrdered(false)

	_, err := c.UserMongoRepository.UploadUsers(ctx, unionID.Hex(), bulkWriteModels, bulkWriteOptions)

	if err != nil {
		return nil, err
	}
	// no need to cache this as well
	return &Response, nil
}

// function to Approve user(member) from client side
func (c *UserController) ApproveUser(ctx context.Context, unionID primitive.ObjectID, memberID primitive.ObjectID) (*model.User, error) {
	if unionID.IsZero() || memberID.IsZero() {
		err := fmt.Errorf("memberID and unionID both are required")
		return nil, err
	}
	//convert unionID to string
	unionIdentifier := unionID.Hex()

	// find member from union's member collection

	member, _ := c.UserMongoRepository.GetMemberByID(ctx, unionIdentifier, memberID)
	// activate the user
	member.Status = "active"

	user, _ := c.UserMongoRepository.Create(ctx, unionIdentifier, member)
	// cache it
	go c.UserRedisRepository.CacheUser(ctx, user.ID.Hex(), user)
	return user, nil
}

func (c *UserController) UpdateUser(ctx context.Context, userID primitive.ObjectID, unionID primitive.ObjectID, update model.UserUpdateInput) (*model.User, error) {
	if userID.IsZero() || unionID.IsZero() {
		err := fmt.Errorf("userID and unionID both are required")
		return nil, err
	}

	updatedUser, _ := c.UserMongoRepository.Update(ctx, unionID.Hex(), userID, update)
	// check if existing cache
	user, _ := c.UserRedisRepository.CacheExists(ctx, userID.Hex())
	if user {
		// invalidate if exists
		go c.UserRedisRepository.InvalidateCache(ctx, userID.Hex())
	}
	return updatedUser, nil
}

func (c *UserController) Login(ctx context.Context, input *model.Credential, device *string) (*model.SingleUserAuth, error) {
	// get the user first
	if input.UnionID.IsZero() {
		err := fmt.Errorf("UnionID is required")
		return nil, err
	}
	unionID := input.UnionID.Hex()

	user, _ := c.UserMongoRepository.GetByUsername(ctx, unionID, input.Username)

	// find the hashedpassword
	hashedPassword, _ := auth.HashPassword(input.Password, unionID)
	// verify the password
	if !auth.VerifyPassword(hashedPassword, input.Password, unionID) {
		err := errors.New("the password is invalid, please try with correct password")
		return nil, err
	}
	// generate token
	token, err := auth.GenerateJWTToken(input.Username, input.UnionID, user.UnionID, 24)
	if err != nil {
		return nil, err
	}

	authenticatedUser := &model.SingleUserAuth{
		User:  *user,
		Token: token,
	}

	return authenticatedUser, nil
}

func (c *UserController) LoginWithToken(ctx context.Context, token *string) (*model.SingleUserAuth, error) {
	userClaim, err := auth.ValidateJWTToken(*token)
	if err != nil {
		err = fmt.Errorf("session expired please login again")
		return nil, err
	}

	user, _ := c.UserMongoRepository.GetByUsername(ctx, userClaim.UnionID.Hex(), userClaim.Username)

	refreshedToken, _ := auth.RefreshJWTToken(*token, 24)
	authenticated := model.SingleUserAuth{
		User:  *user,
		Token: refreshedToken,
	}
	return &authenticated, nil
}

func (c *UserController) ResetPassword(ctx context.Context, unionID primitive.ObjectID, resetKey *string, password *string) (*string, error) {
	filter := bson.M{"passwordResetKey": resetKey}
	//get user by passwordResetKey

	user, err := c.UserMongoRepository.GetUser(ctx, unionID.Hex(), filter)

	if err != nil {
		err = fmt.Errorf("could not find user with passwordResetKey")
		return nil, err
	}
	if user.PasswordResetExpireTime.Before(time.Now()) {
		err = fmt.Errorf("password reset expired")
		return nil, err
	}

	newPassword, err := auth.HashPassword(*password, unionID.Hex())
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"password":                newPassword,
			"passwordResetExpireTime": time.Date(1999, 01, 01, 00, 00, 00, 00, &time.Location{}),
			"passwordResetKey":        "",
			"resetRequired":           false,
		},
	}

	_, err = c.UserMongoRepository.UpdateUser(ctx, unionID.Hex(), filter, update)
	// check if existing cache
	cacheUser, _ := c.UserRedisRepository.CacheExists(ctx, user.ID.Hex())
	if cacheUser {
		// invalidate if exists
		go c.UserRedisRepository.InvalidateCache(ctx, user.ID.Hex())
	}

	if err != nil {
		err = fmt.Errorf("could not update User")
		return nil, err
	}
	return &Response, nil
}

func (c *UserController) ResetPasswordRequest(ctx context.Context, unionID primitive.ObjectID, username string) (*string, error) {
	//verify user existence
	user, err := c.UserMongoRepository.GetByUsername(ctx, unionID.Hex(), username)
	if err != nil {
		err = fmt.Errorf("could not find user with username %v", username)
		return nil, err
	}
	resetPasswordLink := os.Getenv("PWD_RESET_PATH")
	db, err := c.dbManager.GetDatabase(ctx, unionID.Hex())
	if err != nil {
		return nil, err
	}
	unionName := db.Name()
	mailContent := email.GetResetPasswordBody(unionName, username, resetPasswordLink)
	mutationInput := map[string]interface{}{
		"email":    user.Profile.Email,
		"subject":  "Request Password Reset",
		"content":  mailContent,
		"category": "password",
	}

	gqlEP := os.Getenv("Comm_GRAPHQL_ENDPOINT")
	c.graphqlManager.SetgqlEndpoint(gqlEP)
	mailMutationBuilder := c.graphqlManager.GetMutationBuilder()
	// Create a user
	createUserMutation, createVars := mailMutationBuilder.
		SetMutationName("sendMail").
		SetInputName("SendMailInput").
		SetInput(mutationInput).
		Build()
	var result struct {
		Response string `json:"sendMail"`
	}

	err = c.graphqlManager.Execute(ctx, createUserMutation, createVars, &result)
	if err != nil {
		err = fmt.Errorf("could not complete request for password reset please try after some time or contact Admin")
		return nil, err
	}
	return &result.Response, err
}

func (c *UserController) DeleteUser(ctx context.Context, userID primitive.ObjectID, unionID primitive.ObjectID) (string, error) {
	if userID.IsZero() || unionID.IsZero() {
		err := fmt.Errorf("userID and unionID both are equally required")
		return "", err
	}

	err := c.UserMongoRepository.Delete(ctx, unionID.Hex(), userID)

	if err != nil {
		err = fmt.Errorf("failed to delete user due to %v", err)
		return "", err
	}
	go c.UserRedisRepository.InvalidateCache(ctx, userID.Hex())
	return Response, err
}

func (c *UserController) RestoreUser(ctx context.Context, userID primitive.ObjectID, unionID primitive.ObjectID) (string, error) {
	if userID.IsZero() || unionID.IsZero() {
		err := fmt.Errorf("userID and unionID both are equally required")
		return "", err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted": false,
		},
	}

	err := c.UserMongoRepository.Restore(ctx, unionID.Hex(), userID, update)

	if err != nil {
		err = fmt.Errorf("failed to delete user due to %v", err)
		return "", err
	}
	// check if existing cache
	cacheUser, _ := c.UserRedisRepository.CacheExists(ctx, userID.Hex())
	if cacheUser {
		// invalidate if exists
		go c.UserRedisRepository.InvalidateCache(ctx, userID.Hex())
	}
	return Response, nil
}

func (c *UserController) User(ctx context.Context, userID primitive.ObjectID, unionID primitive.ObjectID) (*model.User, error) {

	user, err := c.UserMongoRepository.GetByID(ctx, unionID.Hex(), userID)
	if err != nil {
		return nil, err
	}
	go c.UserRedisRepository.CacheUser(ctx, user.ID.Hex(), user)
	return user, nil
}

func (c *UserController) Users(ctx context.Context, filter *model.UserFilterInput, page int, limit int) ([]*model.User, error) {
	users, err := c.UserMongoRepository.Find(ctx, filter.UnionID.Hex(), filter, page, limit)
	if err != nil {
		return nil, err
	}
	go c.UserRedisRepository.CacheUsers(ctx, "all-users-"+filter.UnionID.Hex(), users) // static key as all-users-unionID
	return users, nil
}

func (r *UserController) UserCount(ctx context.Context, filter *model.UserFilterInput) (int64, error) {
	return r.UserMongoRepository.Count(ctx, string(filter.UnionID.Hex()), filter)

}

// HELPER
func boolPtr(b bool) *bool {
	return &b
}
