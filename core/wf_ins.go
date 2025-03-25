package core

import (
	"context"
	"encoding/json"
	"time"

	"wf_demo/infra/dao"
	"wf_demo/infra/dto"
	"wf_demo/infra/util"
)

type WfInsManager interface {
	CreateWfInst(ctx context.Context, wfDefID int, wfDefName, title string) (int, error)
	GetWfInstInfo(ctx context.Context, wfInstID int) (*dao.WfIns, error)
	GetNodeInfos(ctx context.Context, wfInstID int) ([]*NodeInfo, error)
	GetAllPendingInsts(ctx context.Context) ([]*dao.WfIns, error)
	MoveStage(ctx context.Context, nodInfos []*NodeInfo, taskID, wfInstID, step int, pass bool) error
}

type WfInsManagerImpl struct {
	wfExecDTO        dto.WfExecDTO
	wfInstDTO        dto.WfInsDTO
	wfTaskDTO        dto.WfTaskDTO
	wfParticipantDTO dto.WfParticipantDTO
}

func NewWfInsManager() WfInsManager {
	return &WfInsManagerImpl{
		wfExecDTO:        dto.NewWfExecDTO(),
		wfInstDTO:        dto.NewWfInsDTO(),
		wfTaskDTO:        dto.NewWfTaskDTO(),
		wfParticipantDTO: dto.NewWfParticipantDTO(),
	}
}

func (s *WfInsManagerImpl) CreateWfInst(ctx context.Context, wfDefID int, wfDefName, title string) (int, error) {
	wfIns := &dao.WfIns{
		WfDefId:   wfDefID,
		WfDefName: wfDefName,
		Title:     title,
		StartTime: time.Now(),
	}
	return s.wfInstDTO.CreateWfInst(ctx, wfIns)
}

func (s *WfInsManagerImpl) GetWfInstInfo(ctx context.Context, wfInstID int) (*dao.WfIns, error) {
	wfInst, err := s.wfInstDTO.GetWfInstByID(ctx, wfInstID)
	if err != nil {
		util.GetLogger(ctx).Errorf("GetWfInstByID failed: %v", err)
		return nil, err
	}
	return wfInst, nil
}

func (s *WfInsManagerImpl) GetNodeInfos(ctx context.Context, wfInstID int) ([]*NodeInfo, error) {
	wfExecRecord, err := s.wfExecDTO.GetWfExecByInstID(ctx, wfInstID)
	if err != nil {
		util.GetLogger(ctx).Errorf("GetNodeInfos failed: %v", err)
		return nil, err
	}
	var nodeInfos []*NodeInfo
	err = json.Unmarshal([]byte(wfExecRecord.NodeInfos), &nodeInfos)
	if err != nil {
		util.GetLogger(ctx).Errorf("Unmarshal nodeInfos failed: %v", err)
		return nil, err
	}
	return nodeInfos, nil
}

func (s *WfInsManagerImpl) GetAllPendingInsts(ctx context.Context) ([]*dao.WfIns, error) {
	return s.wfInstDTO.GetPendingWfInstList(ctx)
}

func (s *WfInsManagerImpl) MoveStage(ctx context.Context, nodeInfos []*NodeInfo, taskID, wfInstID, step int, pass bool) error {
	if pass {
		step++
	} else {
		step--
	}
	// 执行下一个节点
	if nodeInfos[step].ApproverType == NOTIFIER {
		// 创建通知任务
		notifyTask := &dao.WfTask{
			NodeID:     nodeInfos[step].NodeID,
			Step:       step,
			WfInstID:   wfInstID,
			IsFinished: true,
		}
		_, err := s.wfTaskDTO.CreateWfTask(ctx, notifyTask)
		if err != nil {
			util.GetLogger(ctx).Errorf("Create notify task failed: %v", err)
			return err
		}
		_, err = s.wfParticipantDTO.AddNotifier(ctx, wfInstID, taskID, step, nodeInfos[step].Approver)
		if err != nil {
			util.GetLogger(ctx).Errorf("Add notifier failed: %v", err)
			return err
		}
		util.GetLogger(ctx).Infof("CompleteTask: %v", notifyTask)
		return s.MoveStage(ctx, nodeInfos, taskID, wfInstID, step, pass)
	}
	if pass {
		// 执行下一个节点
		return s.move2NextStage(ctx, nodeInfos, wfInstID, step)
	}
	return s.move2PrevStage(ctx, nodeInfos, wfInstID, step)
}

func (s *WfInsManagerImpl) move2NextStage(ctx context.Context, nodeInfos []*NodeInfo, wfInstID, step int) error {
	newTask := &dao.WfTask{
		NodeID:     nodeInfos[step].NodeID,
		Step:       step,
		WfInstID:   wfInstID,
		IsFinished: false,
	}
	wfIns := &dao.WfIns{
		NodeID:    nodeInfos[step].NodeID,
		Candidate: nodeInfos[step].Approver,
	}
	wfIns.ID = wfInstID
	if step+1 != len(nodeInfos) {
		// 不是最后一步
		taskID, err := s.wfTaskDTO.CreateWfTask(ctx, newTask)
		if err != nil {
			util.GetLogger(ctx).Errorf("Create task failed: %v", err)
			return err
		}
		wfIns.TaskID = taskID
		err = s.wfInstDTO.UpdateWfInst(ctx, wfIns)
		if err != nil {
			util.GetLogger(ctx).Errorf("Update wfIns failed: %v", err)
			return err
		}
	} else {
		// 最后一步，直接结束
		newTask.IsFinished = true
		taskID, err := s.wfTaskDTO.CreateWfTask(ctx, newTask)
		if err != nil {
			util.GetLogger(ctx).Errorf("Create task failed: %v", err)
			return err
		}
		// 更新wfIns
		wfIns.TaskID = taskID
		wfIns.IsFinished = true
		wfIns.EndTime = time.Now()
		err = s.wfInstDTO.UpdateWfInst(ctx, wfIns)
		if err != nil {
			util.GetLogger(ctx).Errorf("Update wfIns failed: %v", err)
			return err
		}
	}
	return nil
}

func (s *WfInsManagerImpl) move2PrevStage(ctx context.Context, nodeInfos []*NodeInfo, wfInstID, step int) error {
	if step <= 0 {
		util.GetLogger(ctx).Infof("已经回退至流程起点，无法再回退")
		return nil
	}
	newTask := &dao.WfTask{
		NodeID:     nodeInfos[step].NodeID,
		Step:       step,
		WfInstID:   wfInstID,
		IsFinished: false,
	}
	taskID, err := s.wfTaskDTO.CreateWfTask(ctx, newTask)
	if err != nil {
		util.GetLogger(ctx).Errorf("Create task failed: %v", err)
		return err
	}
	wfIns := &dao.WfIns{
		NodeID:    nodeInfos[step].NodeID,
		Candidate: nodeInfos[step].Approver,
		TaskID:    taskID,
	}
	wfIns.ID = wfInstID
	err = s.wfInstDTO.UpdateWfInst(ctx, wfIns)
	if err != nil {
		util.GetLogger(ctx).Errorf("Update wfIns failed: %v", err)
		return err
	}
	return nil
}
