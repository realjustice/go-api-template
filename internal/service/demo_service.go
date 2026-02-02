package service

import (
	"context"

	"go-api-template/internal/model"
	"go-api-template/internal/repository"
	"go-api-template/pkg/errors"
	"go-api-template/pkg/logger"
)

// DemoService Demo 业务逻辑层
type DemoService struct {
	demoRepo *repository.DemoRepository
}

// NewDemoService 创建 Demo Service
func NewDemoService(demoRepo *repository.DemoRepository) *DemoService {
	return &DemoService{
		demoRepo: demoRepo,
	}
}

// GetByID 根据 ID 获取
func (s *DemoService) GetByID(ctx context.Context, id uint) (*model.Demo, error) {
	demo, err := s.demoRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("get demo by id failed",
			logger.Uint("id", id),
			logger.Err(err),
		)
		return nil, err
	}
	return demo, nil
}

// GetAll 获取所有
func (s *DemoService) GetAll(ctx context.Context) ([]*model.Demo, error) {
	demos, err := s.demoRepo.FindAll(ctx)
	if err != nil {
		logger.Error("get all demos failed", logger.Err(err))
		return nil, err
	}
	return demos, nil
}

// Create 创建
func (s *DemoService) Create(ctx context.Context, demo *model.Demo) error {
	// 业务逻辑校验
	if demo.Title == "" {
		return errors.New("title cannot be empty")
	}

	err := s.demoRepo.Create(ctx, demo)
	if err != nil {
		logger.Error("create demo failed",
			logger.String("title", demo.Title),
			logger.Err(err),
		)
		return err
	}

	logger.Info("demo created successfully",
		logger.Uint("id", demo.ID),
		logger.String("title", demo.Title),
	)
	return nil
}

// Update 更新
func (s *DemoService) Update(ctx context.Context, id uint, demo *model.Demo) error {
	// 检查是否存在
	existing, err := s.demoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新字段
	existing.Title = demo.Title
	existing.Content = demo.Content
	existing.Status = demo.Status

	err = s.demoRepo.Update(ctx, existing)
	if err != nil {
		logger.Error("update demo failed",
			logger.Uint("id", id),
			logger.Err(err),
		)
		return err
	}

	logger.Info("demo updated successfully", logger.Uint("id", id))
	return nil
}

// Delete 删除
func (s *DemoService) Delete(ctx context.Context, id uint) error {
	// 检查是否存在
	_, err := s.demoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.demoRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("delete demo failed",
			logger.Uint("id", id),
			logger.Err(err),
		)
		return err
	}

	logger.Info("demo deleted successfully", logger.Uint("id", id))
	return nil
}
