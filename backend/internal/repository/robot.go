package repository

import (
	"backend/internal/model"
	"context"
)

type RobotRepository interface {
	GetRobot(ctx context.Context, id int64) (*model.Robot, error)
	CreateRobot(ctx context.Context, robot *model.Robot) error
	UpdateRobot(ctx context.Context, robot *model.Robot) error
	DeleteRobot(ctx context.Context, id int64) error
	ListRobots(ctx context.Context, page int, size int, options map[string]interface{}) ([]*model.Robot, int64, error)
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
	if err := r.DB(ctx).First(&robot, id).Error; err != nil {
		return nil, err
	}
	return &robot, nil
}

func (r *robotRepository) CreateRobot(ctx context.Context, robot *model.Robot) error {
	err := r.DB(ctx).Create(robot).Error
	return err
}

func (r *robotRepository) UpdateRobot(ctx context.Context, robot *model.Robot) error {
	err := r.DB(ctx).Model(robot).Updates(map[string]interface{}{
		"name":     robot.Name,
		"desc":     robot.Desc,
		"webhook":  robot.Webhook,
		"callback": robot.Callback,
		"options":  robot.Options,
		"enabled":  robot.Enabled,
		"owner":    robot.Owner,
	}).Error
	return err
}

func (r *robotRepository) DeleteRobot(ctx context.Context, id int64) error {
	err := r.DB(ctx).Delete(&model.Robot{}, id).Error
	return err
}

func (r *robotRepository) ListRobots(ctx context.Context, page int, size int, options map[string]interface{}) ([]*model.Robot, int64, error) {
	var (
		robots []*model.Robot
		total  int64
	)

	stmtCount := r.DB(ctx).Model(&model.Robot{})
	if options != nil {
		for key, value := range options {
			stmtCount = stmtCount.Where(key+" = ?", value)
		}
	}
	if err := stmtCount.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	stmtData := r.DB(ctx).Limit(size).Offset(offset)
	if options != nil {
		for key, value := range options {
			stmtData = stmtData.Where(key+" = ?", value)
		}
	}
	if err := stmtData.Find(&robots).Error; err != nil {
		return nil, 0, err
	}

	return robots, total, nil
}
