package dto

import (
	"context"

	"wf_demo/infra/dao"
)

type WfTaskDTO interface {
	CreateWfTask(ctx context.Context, task *dao.WfTask) (int, error)
	GetWfTaskByID(ctx context.Context, taskID int) (*dao.WfTask, error)
	UpdateWfTask(ctx context.Context, task *dao.WfTask) error
}

type WfTaskDTOImpl struct {
}

func NewWfTaskDTO() WfTaskDTO {
	return &WfTaskDTOImpl{}
}

func (w *WfTaskDTOImpl) CreateWfTask(ctx context.Context, task *dao.WfTask) (int, error) {
	return dao.Create(dao.GetTxFromContext(ctx), task)
}

func (w *WfTaskDTOImpl) GetWfTaskByID(ctx context.Context, taskID int) (*dao.WfTask, error) {
	records, err := dao.GetByCond[dao.WfTask](dao.GetTxFromContext(ctx), map[string]any{
		dao.ID: taskID,
	})
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return records[0], nil
}

func (w *WfTaskDTOImpl) UpdateWfTask(ctx context.Context, task *dao.WfTask) error {
	return dao.Update(dao.GetTxFromContext(ctx), task)
}
