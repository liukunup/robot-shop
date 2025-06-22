package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"

	"gorm.io/gorm"
)

type RobotService interface {
	List(ctx context.Context, req *v1.ListRobotRequest) (*v1.ListRobotResponseData, error)
	Create(ctx context.Context, req *v1.RobotCreateRequest) error
	Update(ctx context.Context, req *v1.RobotUpdateRequest) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (model.Robot, error)
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

func (s *robotService) List(ctx context.Context, req *v1.ListRobotRequest) (*v1.ListRobotResponseData, error) {
	list, total, err := s.robotRepository.List(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.ListRobotResponseData{
		List:  make([]v1.RobotDataItem, 0),
		Total: total,
	}
	for _, robot := range list {
		data.List = append(data.List, v1.RobotDataItem{
			Id:        robot.ID,
			Name:      robot.Name,
			Desc:      robot.Desc,
			Webhook:   robot.Webhook,
			Callback:  robot.Callback,
			Enabled:   robot.Enabled,
			CreatedAt: robot.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: robot.UpdatedAt.Format(constant.DateTimeLayout),
		})
	}
	return data, nil
}

func (s *robotService) Create(ctx context.Context, req *v1.RobotCreateRequest) error {
	return s.robotRepository.Create(ctx, &model.Robot{
		Name:     req.Name,
		Desc:     req.Desc,
		Webhook:  req.Webhook,
		Callback: req.Callback,
		Enabled:  req.Enabled,
		Owner:    req.Owner,
	})
}

func (s *robotService) Update(ctx context.Context, req *v1.RobotUpdateRequest) error {
	return s.robotRepository.Update(ctx, &model.Robot{
		Name:     req.Name,
		Desc:     req.Desc,
		Webhook:  req.Webhook,
		Callback: req.Callback,
		Enabled:  req.Enabled,
		Owner:    req.Owner,
		Model: gorm.Model{
			ID: req.ID,
		},
	})
}

func (s *robotService) Delete(ctx context.Context, id uint) error {
	return s.robotRepository.Delete(ctx, id)
}

func (s *robotService) Get(ctx context.Context, id uint) (model.Robot, error) {
	return s.robotRepository.Get(ctx, id)
}
