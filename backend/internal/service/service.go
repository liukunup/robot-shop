package service

import (
	"backend/internal/repository"
	"backend/pkg/email"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"backend/pkg/sid"
	"backend/pkg/storage"
)

type Service struct {
	logger  *log.Logger
	sid     *sid.Sid
	jwt     *jwt.JWT
	email   *email.Email
	storage *storage.Storage
	tm      repository.Transaction
}

func NewService(
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
	email *email.Email,
	storage *storage.Storage,
	tm repository.Transaction,
) *Service {
	return &Service{
		logger:  logger,
		sid:     sid,
		jwt:     jwt,
		email:   email,
		storage: storage,
		tm:      tm,
	}
}
