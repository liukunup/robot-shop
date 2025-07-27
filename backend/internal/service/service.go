package service

import (
	"backend/internal/repository"
	"backend/pkg/email"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"backend/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	email  *email.Email
	tm     repository.Transaction
}

func NewService(
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
	email *email.Email,
	tm repository.Transaction,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		email:  email,
		tm:     tm,
	}
}
