package dto

import (
	"context"
	"time"

	"wf_demo/infra/dao"
)

type WfInsDTO interface {
	CreateWfInst(ctx context.Context, data *dao.WfIns) (int, error)
	UpdateWfInst(ctx context.Context, data *dao.WfIns) error
	GetWfInstByID(ctx context.Context, wfInstID int) (*dao.WfIns, error)
	GetPendingWfInstList(ctx context.Context) ([]*dao.WfIns, error)
}

type WfInsDTOImpl struct {
}

func NewWfInsDTO() WfInsDTO {
	return &WfInsDTOImpl{}
}

func (w *WfInsDTOImpl) CreateWfInst(ctx context.Context, data *dao.WfIns) (int, error) {
	data.StartTime = time.Now()
	return dao.Create(dao.GetTxFromContext(ctx), data)
}

func (w *WfInsDTOImpl) UpdateWfInst(ctx context.Context, data *dao.WfIns) error {
	return dao.Update(dao.GetTxFromContext(ctx), data)
}

func (w *WfInsDTOImpl) GetWfInstByID(ctx context.Context, wfInstID int) (*dao.WfIns, error) {
	records, err := dao.GetByCond[dao.WfIns](dao.GetTxFromContext(ctx), map[string]any{
		dao.ID: wfInstID,
	})
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return records[0], nil
}

func (w *WfInsDTOImpl) GetPendingWfInstList(ctx context.Context) ([]*dao.WfIns, error) {
	return dao.GetByCond[dao.WfIns](dao.GetTxFromContext(ctx), map[string]interface{}{
		"is_finished": false,
	})
}
