package repository

import (
    "context"
	"backend/internal/model"
)

type MessageRepository interface {
	GetMessage(ctx context.Context, id int64) (*model.Message, error)
}

func NewMessageRepository(
	repository *Repository,
) MessageRepository {
	return &messageRepository{
		Repository: repository,
	}
}

type messageRepository struct {
	*Repository
}

func (r *messageRepository) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
	var message model.Message

	return &message, nil
}
