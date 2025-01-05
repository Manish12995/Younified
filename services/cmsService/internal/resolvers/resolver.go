package resolver

import (
	"younified-backend/providers/database"
	"younified-backend/services/cmsService/internal/controller"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DBManager     *database.DBManager
	CMSController *controller.CmsController
}

func NewResolver(dbManager *database.DBManager, CMSController *controller.CmsController) *Resolver {
	return &Resolver{
		DBManager:     dbManager,
		CMSController: CMSController,
	}
}
