package service

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/sid"
	"context"
	"encoding/json"
	"time"
)

type RobotService interface {
	GetRobot(ctx context.Context, id int64) (*model.Robot, error)
	CreateRobot(ctx context.Context, req *v1.RobotRequest) (*model.Robot, error)
	UpdateRobot(ctx context.Context, id int64, req *v1.RobotRequest) (*model.Robot, error)
	DeleteRobot(ctx context.Context, id int64) error
	ListRobots(ctx context.Context, page int, size int, options map[string]interface{}) ([]*model.Robot, int64, error)
}

func NewRobotService(
	service *Service,
	robotRepository repository.RobotRepository,
) RobotService {
	return &robotService{
		Service:         service,
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

func (s *robotService) CreateRobot(ctx context.Context, req *v1.RobotRequest) (*model.Robot, error) {
	robotId, err := sid.NewSid().GenString()
	if err != nil {
		return nil, err
	}

	options, err := json.Marshal(req.Options)
	if err != nil {
		return nil, err
	}

	robot := &model.Robot{
		RobotId:   robotId,
		Name:      req.Name,
		Desc:      req.Desc,
		Webhook:   req.Webhook,
		Callback:  req.Callback,
		Options:   string(options),
		Enabled:   req.Enabled,
		Owner:     req.Owner,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.robotRepository.CreateRobot(ctx, robot); err != nil {
		return nil, err
	}
	return robot, nil
}

func (s *robotService) UpdateRobot(ctx context.Context, id int64, req *v1.RobotRequest) (*model.Robot, error) {
	robot, err := s.robotRepository.GetRobot(ctx, id)
	if err != nil {
		return nil, err
	}

	options, err := json.Marshal(req.Options)
	if err != nil {
		return nil, err
	}

	robot.Name = req.Name
	robot.Desc = req.Desc
	robot.Webhook = req.Webhook
	robot.Callback = req.Callback
	robot.Options = string(options)
	robot.Enabled = req.Enabled
	robot.Owner = req.Owner
	robot.UpdatedAt = time.Now()

	if err := s.robotRepository.UpdateRobot(ctx, robot); err != nil {
		return nil, err
	}
	return robot, nil
}

func (s *robotService) DeleteRobot(ctx context.Context, id int64) error {
	return s.robotRepository.DeleteRobot(ctx, id)
}

func (s *robotService) ListRobots(ctx context.Context, page int, size int, options map[string]interface{}) ([]*model.Robot, int64, error) {
	robots, total, err := s.robotRepository.ListRobots(ctx, page, size, options)
	if err != nil {
		return nil, 0, err
	}
	return robots, total, nil
}
