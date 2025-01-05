package resolver

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"younified-backend/providers/database"
	controller "younified-backend/services/unionService/internal/controller"
)

type Resolver struct {
	DBManager       *database.DBManager
	UnionController *controller.UnionController
}

func NewResolver(dbManager *database.DBManager, UnionController *controller.UnionController) *Resolver {
	return &Resolver{
		DBManager:       dbManager,
		UnionController: UnionController,
	}
}
