package service

import (
    "context"
	"backend/internal/model"
	"backend/internal/repository"
)

type MessageService interface {
	GetMessage(ctx context.Context, id int64) (*model.Message, error)
}
func NewMessageService(
    service *Service,
    messageRepository repository.MessageRepository,
) MessageService {
	return &messageService{
		Service:        service,
		messageRepository: messageRepository,
	}
}

type messageService struct {
	*Service
	messageRepository repository.MessageRepository
}

func (s *messageService) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
	return s.messageRepository.GetMessage(ctx, id)
}
