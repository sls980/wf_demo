package core

import (
	"context"

	"wf_demo/infra/dao"
	"wf_demo/infra/dto"
	"wf_demo/infra/util"
)

type WfTaskService interface {
	CreateWfTask(ctx context.Context, task *dao.WfTask) (int, error)
	GetWfTaskByID(ctx context.Context, taskID int) (*dao.WfTask, error)
	CompleteTask(ctx context.Context, taskID int) error
}

type WfTaskServiceImpl struct {
	wfTaskDTO dto.WfTaskDTO
}

func NewWfTaskServiceImpl() WfTaskService {
	return &WfTaskServiceImpl{
		wfTaskDTO: dto.NewWfTaskDTO(),
	}
}

func (s *WfTaskServiceImpl) CreateWfTask(ctx context.Context, task *dao.WfTask) (int, error) {
	return s.wfTaskDTO.CreateWfTask(ctx, task)
}

func (s *WfTaskServiceImpl) GetWfTaskByID(ctx context.Context, taskID int) (*dao.WfTask, error) {
	return s.wfTaskDTO.GetWfTaskByID(ctx, taskID)
}

func (s *WfTaskServiceImpl) CompleteTask(ctx context.Context, taskID int) error {
	taskInst, err := s.GetWfTaskByID(ctx, taskID)
	if err != nil {
		util.GetLogger(ctx).Errorf("GetWfTaskByID failed: %v", err)
		return err
	}
	taskInst.IsFinished = true
	util.GetLogger(ctx).Infof("CompleteTask: %v", taskInst)
	return s.wfTaskDTO.UpdateWfTask(ctx, taskInst)
}
