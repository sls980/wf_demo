package dto

import (
	"context"

	"wf_demo/infra/dao"
)

type WfParticipantDTO interface {
	AddParticipant(ctx context.Context, wfInsID, taskID, step int, userID string) (int, error)
	AddNotifier(ctx context.Context, wfInsID, taskID, step int, userID string) (int, error)
	AddCandidate(ctx context.Context, wfInsID, taskID, step int, userID string) (int, error)
}

type WfParticipantDTOImpl struct {
}

func NewWfParticipantDTO() WfParticipantDTO {
	return &WfParticipantDTOImpl{}
}

func (w *WfParticipantDTOImpl) AddParticipant(ctx context.Context, wfInsID, taskID, step int, userID string) (int, error) {
	data := &dao.WfParticipant{
		Type:    dao.PT_Participant,
		UserID:  userID,
		TaskID:  taskID,
		Step:    step,
		WfInsID: wfInsID,
	}
	return dao.Create(dao.GetTxFromContext(ctx), data)
}

func (w *WfParticipantDTOImpl) AddNotifier(ctx context.Context, wfInsID, taskID, step int, userID string) (int, error) {
	data := &dao.WfParticipant{
		Type:    dao.PT_Notifier,
		UserID:  userID,
		TaskID:  taskID,
		Step:    step,
		WfInsID: wfInsID,
	}
	return dao.Create(dao.GetTxFromContext(ctx), data)
}

func (w *WfParticipantDTOImpl) AddCandidate(ctx context.Context, wfInsID, taskID, step int, userID string) (int, error) {
	data := &dao.WfParticipant{
		Type:    dao.PT_Candidate,
		UserID:  userID,
		TaskID:  taskID,
		Step:    step,
		WfInsID: wfInsID,
	}
	return dao.Create(dao.GetTxFromContext(ctx), data)
}
