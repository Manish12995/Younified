package resolver

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"younified-backend/providers/database"
	controller "younified-backend/services/userService/internal/controller"
)

type Resolver struct {
	DBManager      *database.DBManager
	UserController *controller.UserController
}

func NewResolver(dbManager *database.DBManager, UserController *controller.UserController) *Resolver {
	return &Resolver{
		DBManager:      dbManager,
		UserController: UserController,
	}
}
