package session

import (
	"fmt"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/s3-explorer/backend/src/utils"
)

const ServiceName = "session"

type IService interface {
	gousu.IService

	GetSessionToken() string
	VerifySessionToken(sessionToken string) error
}

type Service struct {
	log          *logger.Log
	sessionToken string
}

var _ IService = (*Service)(nil)

func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) Start() error {
	var err error

	s.sessionToken, err = utils.RandString(128)
	if err != nil {
		return fmt.Errorf("error generating session token: %s", err)
	}

	return nil
}

// Health returns always nil (healthy)
func (s *Service) Health() error {
	return nil
}

// Stop does nothing
func (s *Service) Stop() error {
	return nil
}

// NewService creates a new instance of the session Service
func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ (gousu.ServiceFactory) = NewService
