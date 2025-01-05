package controller

import (
	"context"
	"fmt"
	"younified-backend/contracts/communication/model"
	"younified-backend/providers/database"
	"younified-backend/providers/graphqlclient"

	// "younified-backend/providers/redis"
	"younified-backend/services/communicationService/internal/repository"
)

var Response string = "Operation Successful"

type CommsController struct {
	CommsRepository *repository.SendgridRepository

	dbManager      *database.DBManager
	graphqlManager *graphqlclient.Graph
	// redisProvider   *redis.Provider
}

func NewCommsController(dbManager *database.DBManager) *CommsController {
	if dbManager == nil {
		panic("dbManager cannot be nil")
	}

	return &CommsController{
		CommsRepository: repository.NewSendgridRepository(),
		dbManager:       dbManager,
		// redisProvider:   redisProvider,
	}
}

func (c *CommsController) SendMail(ctx context.Context, mailInput model.SendMail) (string, error) {
	var receiver []string
	receiver = append(receiver, mailInput.Email)
	opts := repository.EmailOptions{
		To:      receiver,
		HTML:    mailInput.Content,
		Subject: mailInput.Subject,
	}
	err := c.CommsRepository.SendEmail(opts)
	if err != nil {
		err = fmt.Errorf("we are unable to complete the request to reset password please try after some time")
		return "", err
	}
	return Response, nil
}
