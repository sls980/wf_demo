package core

import (
	"context"
	"encoding/json"

	"wf_demo/infra/dao"
	"wf_demo/infra/dto"
	"wf_demo/infra/util"
)

type WfExecService interface {
	CreateWfFlowList(ctx context.Context, wfDefID, wfInstID int, wfDefName, wfDefSetting string, variables map[string]string) ([]*NodeInfo, error)
}

type WfExecServiceImpl struct {
	wfDefParser WfDefParser
	wfExecDTO   dto.WfExecDTO
}

func NewWfExecServiceImpl() WfExecService {
	return &WfExecServiceImpl{
		wfDefParser: NewWfDefParser(),
		wfExecDTO:   dto.NewWfExecDTO(),
	}
}

func (s *WfExecServiceImpl) CreateWfFlowList(ctx context.Context, wfDefID, wfInstID int, wfDefName, wfDefSetting string, variables map[string]string) ([]*NodeInfo, error) {
	nodeInfos, err := s.wfDefParser.ParseWfDef(ctx, wfDefSetting, variables)
	if err != nil {
		util.GetLogger(ctx).Errorf("解析流程定义失败: %v", err)
		return nil, err
	}
	// 落库
	flowStr, err := json.Marshal(nodeInfos)
	if err != nil {
		util.GetLogger(ctx).Errorf("decode nodeinfos failed")
		return nil, err
	}
	_, err = s.wfExecDTO.CreateWfExec(ctx, &dao.WfExec{
		WfInsId:   wfInstID,
		WfDefId:   wfDefID,
		WfDefName: wfDefName,
		NodeInfos: string(flowStr),
	})
	return nodeInfos, err
}
