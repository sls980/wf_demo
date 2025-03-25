package dto

import (
	"context"
	"time"

	"wf_demo/infra/dao"
)

type WfExecDTO interface {
	CreateWfExec(ctx context.Context, data *dao.WfExec) (int, error)
	GetWfExecByInstID(ctx context.Context, wfInstID int) (*dao.WfExec, error)
}

type WfExecDTOImpl struct {
}

func NewWfExecDTO() WfExecDTO {
	return &WfExecDTOImpl{}
}

func (w *WfExecDTOImpl) CreateWfExec(ctx context.Context, data *dao.WfExec) (int, error) {
	data.IsActive = true
	data.StartTime = time.Now()
	return dao.Create(dao.GetTxFromContext(ctx), data)
}

func (w *WfExecDTOImpl) GetWfExecByInstID(ctx context.Context, wfInstID int) (*dao.WfExec, error) {
	records, err := dao.GetByCond[dao.WfExec](dao.GetTxFromContext(ctx), map[string]any{"wf_ins_id": wfInstID})
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return records[0], nil
}
