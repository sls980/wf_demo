package app

import (
	"context"
	"time"

	"wf_demo/infra/dto"
	"wf_demo/infra/util"
)

// 流程定义服务

type WfDef struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Setting    string     `json:"setting"`
	CreateTime *time.Time `json:"create_time"`
	UpdateTime *time.Time `json:"update_time"`
}

type WfDefService interface {
	// 保存流程定义
	SaveWfDef(ctx context.Context, data *WfDef) (int, error)
	// 根据流程定义id获取流程定义
	GetWfDefById(ctx context.Context, wfDefId int) (*WfDef, error)
	// 查询所有流程定义
	GetAllWfDef(ctx context.Context) ([]*WfDef, error)
	// 删除流程定义
	DeleteWfDef(ctx context.Context, wfDefId int) error
}

type WfDefServiceImpl struct {
	wfDefDTO dto.WfDefDTO
}

func NewWfDefServiceImpl() WfDefService {
	return &WfDefServiceImpl{
		wfDefDTO: dto.NewWfDefDTO(),
	}
}

func (s *WfDefServiceImpl) SaveWfDef(ctx context.Context, data *WfDef) (int, error) {
	wfDefID, err := s.wfDefDTO.CreateWfDef(ctx, data.Name, data.Setting)
	if err != nil {
		util.GetLogger(ctx).Errorf("create error: %v", err)
		return 0, err
	}
	return wfDefID, nil
}

func (s *WfDefServiceImpl) GetWfDefById(ctx context.Context, wfDefId int) (*WfDef, error) {
	wfDef, err := s.wfDefDTO.GetWfDefById(ctx, wfDefId)
	if err != nil {
		util.GetLogger(ctx).Errorf("get wf def by id error: %v", err)
		return nil, err
	}
	var res WfDef
	err = util.Decode(wfDef, &res)
	if err != nil {
		util.GetLogger(ctx).Errorf("decode error: %v", err)
		return nil, err
	}
	return &res, nil
}

func (s *WfDefServiceImpl) GetAllWfDef(ctx context.Context) ([]*WfDef, error) {
	wfDefs, err := s.wfDefDTO.GetAllWfDef(ctx)
	if err != nil {
		util.GetLogger(ctx).Errorf("get all wf def error: %v", err)
		return nil, err
	}
	var res []*WfDef
	for _, wfDef := range wfDefs {
		var wfDefRes WfDef
		err = util.Decode(wfDef, &wfDefRes)
		if err != nil {
			util.GetLogger(ctx).Errorf("decode error: %v", err)
			return nil, err
		}
		res = append(res, &wfDefRes)
	}
	return res, nil
}

func (s *WfDefServiceImpl) DeleteWfDef(ctx context.Context, wfDefId int) error {
	return s.wfDefDTO.DeleteWfDef(ctx, wfDefId)
}
