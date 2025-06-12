package service

import (
    "context"
	"backend/internal/model"
	"backend/internal/repository"
)

type RobotService interface {
	GetRobot(ctx context.Context, id int64) (*model.Robot, error)
}
func NewRobotService(
    service *Service,
    robotRepository repository.RobotRepository,
) RobotService {
	return &robotService{
		Service:        service,
		robotRepository: robotRepository,
	}
}

type robotService struct {
	*Service
	robotRepository repository.RobotRepository
}

func (s *robotService) GetRobot(ctx context.Context, id int64) (*model.Robot, error) {
	return s.robotRepository.GetRobot(ctx, id)
}
