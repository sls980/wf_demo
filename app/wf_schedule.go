package app

import (
	"context"

	"wf_demo/core"
	"wf_demo/infra/dao"
	"wf_demo/infra/dto"
	"wf_demo/infra/util"

	"github.com/jinzhu/gorm"
)

// 流程调度器

type StartProcessReq struct {
	WfDefId int
	Params  map[string]string
}

type WfScheduleService interface {
	// 启动流程
	StartProcess(ctx context.Context, req *StartProcessReq) (int, error)
	// 审批流程
	CompleteTask(ctx context.Context, wfInsID int, pass bool) error
	// 获取所有pending的流程
	GetAllPendingInst(ctx context.Context) ([]*dao.WfIns, error)
}

type wfScheduleServiceImpl struct {
	wfDefService     WfDefService
	wfDefParser      core.WfDefParser
	wfInsManager     core.WfInsManager
	wfTaskSrv        core.WfTaskService
	wfExecSrv        core.WfExecService
	wfParticipantDTO dto.WfParticipantDTO
}

func NewWfScheduleServiceImpl() WfScheduleService {
	return &wfScheduleServiceImpl{
		wfDefService:     NewWfDefServiceImpl(),
		wfDefParser:      core.NewWfDefParser(),
		wfInsManager:     core.NewWfInsManager(),
		wfTaskSrv:        core.NewWfTaskServiceImpl(),
		wfExecSrv:        core.NewWfExecServiceImpl(),
		wfParticipantDTO: dto.NewWfParticipantDTO(),
	}
}

func (s *wfScheduleServiceImpl) StartProcess(ctx context.Context, req *StartProcessReq) (int, error) {
	// 获取流程定义
	wfDef, err := s.wfDefService.GetWfDefById(ctx, req.WfDefId)
	if err != nil {
		util.GetLogger(ctx).Errorf("获取流程定义失败: %v", err)
		return 0, err
	}
	// 开启事务
	var wfInsID int
	if err = dao.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// 创建流程实例&落库
		wfInsID, err = s.wfInsManager.CreateWfInst(ctx, wfDef.Id, wfDef.Name, req.Params["title"])
		if err != nil {
			util.GetLogger(ctx).Errorf("创建流程实例失败: %v", err)
			return err
		}
		util.GetLogger(ctx).Infof("创建流程实例成功: %v", wfInsID)
		// 解析流程定义&生成执行流
		nodeInfoList, err := s.wfExecSrv.CreateWfFlowList(ctx, wfDef.Id, wfInsID, wfDef.Name, wfDef.Setting, req.Params)
		if err != nil {
			util.GetLogger(ctx).Errorf("解析流程定义失败: %v", err)
			return err
		}
		// 生成新任务
		taskID, err := s.wfTaskSrv.CreateWfTask(ctx, &dao.WfTask{
			NodeID:     "start",
			Step:       0,
			WfInstID:   wfInsID,
			IsFinished: true,
		})
		if err != nil {
			util.GetLogger(ctx).Errorf("生成任务失败: %v", err)
			return err
		}
		// 添加流程发起人
		_, err = s.wfParticipantDTO.AddParticipant(ctx, wfInsID, taskID, 0, req.Params["user_id"])
		if err != nil {
			util.GetLogger(ctx).Errorf("AddParticipant failed: %v", err)
			return err
		}
		// 流转执行流
		return s.wfInsManager.MoveStage(ctx, nodeInfoList, taskID, wfInsID, 0, true)
	}); err != nil {
		util.GetLogger(ctx).Errorf("事务执行失败: %v", err)
		return 0, err
	}
	return wfInsID, nil
}

func (s *wfScheduleServiceImpl) CompleteTask(ctx context.Context, wfInsID int, pass bool) error {
	if err := dao.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// 查询流程实例详情
		wfInst, err := s.wfInsManager.GetWfInstInfo(ctx, wfInsID)
		if err != nil {
			util.GetLogger(ctx).Errorf("GetWfInstInfo failed: %v", err)
			return err
		}
		if wfInst.IsFinished {
			util.GetLogger(ctx).Errorf("流程已完成")
			return nil
		}
		taskInfo, err := s.wfTaskSrv.GetWfTaskByID(ctx, wfInst.TaskID)
		if err != nil {
			util.GetLogger(ctx).Errorf("GetWfTaskByID failed: %v", err)
			return err
		}
		// 完成任务
		if pass {
			err = s.wfTaskSrv.CompleteTask(ctx, wfInst.TaskID)
			if err != nil {
				util.GetLogger(ctx).Errorf("CompleteTask failed: %v", err)
				return err
			}
			_, err = s.wfParticipantDTO.AddCandidate(ctx, wfInst.ID, wfInst.TaskID, taskInfo.Step, wfInst.Candidate)
			if err != nil {
				util.GetLogger(ctx).Errorf("Add candidate failed: %v", err)
				return err
			}
		}
		// 流转流程
		nodeInfos, err := s.wfInsManager.GetNodeInfos(ctx, wfInsID)
		if err != nil {
			util.GetLogger(ctx).Errorf("GetNodeInfos failed: %v", err)
			return err
		}
		return s.wfInsManager.MoveStage(ctx, nodeInfos, wfInst.TaskID, wfInsID, taskInfo.Step, pass)
	}); err != nil {
		util.GetLogger(ctx).Errorf("开启事务失败: %v", err)
		return err
	}
	return nil
}

func (s *wfScheduleServiceImpl) GetAllPendingInst(ctx context.Context) ([]*dao.WfIns, error) {
	return s.wfInsManager.GetAllPendingInsts(ctx)
}
