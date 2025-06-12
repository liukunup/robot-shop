package repository

import (
    "context"
	"backend/internal/model"
)

type RobotRepository interface {
	GetRobot(ctx context.Context, id int64) (*model.Robot, error)
}

func NewRobotRepository(
	repository *Repository,
) RobotRepository {
	return &robotRepository{
		Repository: repository,
	}
}

type robotRepository struct {
	*Repository
}

func (r *robotRepository) GetRobot(ctx context.Context, id int64) (*model.Robot, error) {
	var robot model.Robot

	return &robot, nil
}
