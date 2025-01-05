package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"younified-backend/contracts/union/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUnion is the resolver for the createUnion field.
func (r *mutationResolver) CreateUnion(ctx context.Context, input model.RegisterInput) (response *model.Union, err error) {
	return r.UnionController.Register(ctx, input)
}

// ModifyUnion is the resolver for the modifyUnion field.
func (r *mutationResolver) ModifyUnion(ctx context.Context, id primitive.ObjectID, union model.Union) (*model.Union, error) {
	return r.UnionController.ModifyUnion(ctx, id, union)
}

// DeleteUnion is the resolver for the deleteUnion field.
func (r *mutationResolver) DeleteUnion(ctx context.Context, id primitive.ObjectID) (*bool, error) {
	return r.UnionController.DeleteUnion(ctx, id)
}

// UnionByID is the resolver for the unionById field.
func (r *queryResolver) UnionByID(ctx context.Context, id primitive.ObjectID) (*model.Union, error) {
	return r.UnionController.UnionByID(ctx, id)
}

// UnionByName is the resolver for the unionByName field.
func (r *queryResolver) UnionByName(ctx context.Context, name string) (*model.Union, error) {
	return r.UnionController.UnionByName(ctx, name)
}

// Unions is the resolver for the unions field.
func (r *queryResolver) Unions(ctx context.Context, page int, limit int) (*model.UnionsResponse, error) {
	return r.UnionController.Unions(ctx, page, limit)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }


