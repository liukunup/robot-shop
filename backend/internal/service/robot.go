package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
)

type RobotService interface {
	ListRobots(ctx context.Context, req *v1.RobotSearchRequest) (*v1.RobotSearchResponseData, error)
	RobotCreate(ctx context.Context, req *v1.RobotRequest) error
	RobotUpdate(ctx context.Context, id uint, req *v1.RobotRequest) error
	RobotDelete(ctx context.Context, id uint) error
	GetRobot(ctx context.Context, id uint) (model.Robot, error)
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

func (s *robotService) ListRobots(ctx context.Context, req *v1.RobotSearchRequest) (*v1.RobotSearchResponseData, error) {
	list, total, err := s.robotRepository.ListRobots(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.RobotSearchResponseData{
		List:  make([]v1.RobotDataItem, 0),
		Total: total,
	}
	for _, robot := range list {
		data.List = append(data.List, v1.RobotDataItem{
			Id:        robot.ID,
			CreatedAt: robot.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: robot.UpdatedAt.Format(constant.DateTimeLayout),
			Name:      robot.Name,
			Desc:      robot.Desc,
			Webhook:   robot.Webhook,
			Callback:  robot.Callback,
			Enabled:   robot.Enabled,
			Owner:     robot.Owner,
		})
	}
	return data, nil
}

func (s *robotService) RobotCreate(ctx context.Context, req *v1.RobotRequest) error {
	return s.robotRepository.RobotCreate(ctx, &model.Robot{
		Name:     req.Name,
		Desc:     req.Desc,
		Webhook:  req.Webhook,
		Callback: req.Callback,
		Enabled:  req.Enabled,
		Owner:    req.Owner,
	})
}

func (s *robotService) RobotUpdate(ctx context.Context, id uint, req *v1.RobotRequest) error {
	data := map[string]interface{}{
		"name":     req.Name,
		"desc":     req.Desc,
		"webhook":  req.Webhook,
		"callback": req.Callback,
		"enabled":  req.Enabled,
		"owner":    req.Owner,
	}
	return s.robotRepository.RobotUpdate(ctx, id, data)
}

func (s *robotService) RobotDelete(ctx context.Context, id uint) error {
	return s.robotRepository.RobotDelete(ctx, id)
}

func (s *robotService) GetRobot(ctx context.Context, id uint) (model.Robot, error) {
	return s.robotRepository.GetRobot(ctx, id)
}
