package resolver

import (
	"younified-backend/providers/database"
	"younified-backend/services/communicationService/internal/controller"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DBManager       *database.DBManager
	CommsController *controller.CommsController
}
