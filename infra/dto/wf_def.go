package dto

import (
	"context"
	"fmt"

	"wf_demo/infra/dao"
)

type WfDefDTO interface {
	CreateWfDef(ctx context.Context, name, setting string) (int, error)
	GetWfDefById(ctx context.Context, wfDefId int) (*dao.WfDef, error)
	GetAllWfDef(ctx context.Context) ([]*dao.WfDef, error)
	DeleteWfDef(ctx context.Context, wfDefId int) error
}

type WfDefDTOImpl struct {
}

func NewWfDefDTO() WfDefDTO {
	return &WfDefDTOImpl{}
}

func (w *WfDefDTOImpl) CreateWfDef(ctx context.Context, name, setting string) (int, error) {
	// 创建流程定义
	wfDefItem := &dao.WfDef{
		Name:    name,
		Setting: setting,
	}
	return dao.Create(dao.GetTxFromContext(ctx), wfDefItem)
}

func (w *WfDefDTOImpl) GetWfDefById(ctx context.Context, wfDefId int) (*dao.WfDef, error) {
	wfRecords, err := dao.GetByCond[dao.WfDef](dao.GetTxFromContext(ctx), map[string]any{
		dao.ID: wfDefId,
	})
	if err != nil {
		return nil, err
	}
	if len(wfRecords) == 0 {
		return nil, fmt.Errorf("wf def %d not exist", wfDefId)
	}
	return wfRecords[0], nil
}

func (w *WfDefDTOImpl) GetAllWfDef(ctx context.Context) ([]*dao.WfDef, error) {
	return dao.GetByCond[dao.WfDef](dao.GetTxFromContext(ctx), nil)
}

func (w *WfDefDTOImpl) DeleteWfDef(ctx context.Context, wfDefId int) error {
	return dao.Delete[dao.WfDef](dao.GetTxFromContext(ctx), wfDefId)
}
