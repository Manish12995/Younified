package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the methods for interacting with User data
type UserRepository interface {
	// Create/Register operations
	Create(ctx context.Context, unionID string, user *User) (*User, error)
	CreateMember(ctx context.Context, unionID string, user *User) (*User, error)

	// Approve operations
	ApproveUser(ctx context.Context, user *User) (*User, error)

	// Read operations
	GetByID(ctx context.Context, unionID string, userID primitive.ObjectID) (*User, error)
	GetByUsername(ctx context.Context, unionID string, username string) (*User, error)
	GetMemberByID(ctx context.Context, unionID string, id primitive.ObjectID) (*User, error)
	UserByEmployeeID(ctx context.Context, employeeID string) (*User, error)
	Find(ctx context.Context, unionID string, filter *UserFilterInput, page, limit int) ([]*User, error)
	Count(ctx context.Context, unionID string, filter *UserFilterInput) (int64, error)

	// Update operations
	Update(ctx context.Context, unionID string, id primitive.ObjectID, updates bson.M) (*User, error)
	UpdatePassword(ctx context.Context, userID primitive.ObjectID, newPassword string) error
	UpdateProfile(ctx context.Context, userID primitive.ObjectID, profile *UserInfo) error
	UpdateStatus(ctx context.Context, userID primitive.ObjectID, status string) error
	UpdateLoginStatus(ctx context.Context, userID primitive.ObjectID, loggedIn bool) error
	UpdatePermissions(ctx context.Context, userID primitive.ObjectID, permissions map[string]int64) error
	UpdateUnionPosition(ctx context.Context, userID primitive.ObjectID, position string) error

	// Delete operations
	Delete(ctx context.Context, unionID string, id primitive.ObjectID) error
	SoftDeleteUser(ctx context.Context, userID primitive.ObjectID) error

	// Authentication related
	Login(ctx context.Context, credential *Credential) (*SingleUserAuth, error)
	GeneratePasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetKey, newPassword string) error

	// Token management
	AddToken(ctx context.Context, userID primitive.ObjectID, token string) error
	RemoveToken(ctx context.Context, userID primitive.ObjectID, token string) error
	LoginWithToken(ctx context.Context, token string) (*User, error)

	// Bulk operations
	BulkCreateUsers(ctx context.Context, users []*User) (*UserUploadReport, error)
	BulkUpdateUsers(ctx context.Context, users []*User) (*UserUploadReport, error)

	// Union-specific operations
	GetUnionUsers(ctx context.Context, unionID primitive.ObjectID, filter *UserFilterInput, page, limit int) ([]*User, int64, error)
}

// Add more methods as required by your service.
