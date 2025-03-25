package dao

import "time"

// 流程定义表
type WfDef struct {
	BaseModel
	Name    string `json:"name" gorm:"column:name"`
	Setting string `json:"setting" gorm:"column:setting"`
}

func (wd *WfDef) TableName() string {
	return "wf_def"
}

// 流程实例执行流表
type WfExec struct {
	BaseModel
	WfInsId   int       `gorm:"column:wf_ins_id" json:"wfInsId"`     // 流程实例id
	WfDefId   int       `gorm:"column:wf_def_id" json:"wfDefId"`     // 流程定义id
	WfDefName string    `gorm:"column:wf_def_name" json:"wfDefName"` // 流程定义名称
	NodeInfos string    `gorm:"column:node_infos" json:"nodeInfos"`  // 节点信息
	IsActive  bool      `gorm:"column:is_active" json:"isActive"`    // 是否激活
	StartTime time.Time `gorm:"column:start_time" json:"startTime"`  // 开始时间
}

func (w *WfExec) TableName() string {
	return "wf_exec"
}

type WfIns struct {
	BaseModel
	WfDefId    int       `json:"wf_def_id" gorm:"column:wf_def_id"`     // 流程定义ID
	WfDefName  string    `json:"wf_def_name" gorm:"column:wf_def_name"` // 流程定义名称
	Title      string    `json:"title" gorm:"column:title"`             // 流程实例标题
	NodeID     string    `json:"node_id" gorm:"column:node_id"`         // 当前节点ID
	Candidate  string    `json:"candidate" gorm:"column:candidate"`     // 当前审批人
	TaskID     int       `json:"task_id" gorm:"column:task_id"`         // 当前任务ID
	IsFinished bool      `json:"is_finished" gorm:"column:is_finished"` // 是否完成
	StartTime  time.Time `json:"start_time" gorm:"column:start_time"`   // 开始时间
	EndTime    time.Time `json:"end_time" gorm:"column:end_time"`       // 结束时间
}

func (w *WfIns) TableName() string {
	return "wf_ins"
}

type WfParticipant struct {
	BaseModel
	Type    string `json:"type" gorm:"column:type"`
	UserID  string `json:"user_id" gorm:"column:user_id"`
	TaskID  int    `json:"task_id" gorm:"column:task_id"`
	Step    int    `json:"step" gorm:"column:step"`
	WfInsID int    `json:"wf_ins_id" gorm:"column:wf_ins_id"`
}

func (wp *WfParticipant) TableName() string {
	return "wf_participant"
}

const (
	PT_Candidate   = "candidate"
	PT_Participant = "participant"
	PT_Notifier    = "notifier"
)

type WfTask struct {
	BaseModel
	NodeID     string `json:"node_id" gorm:"column:node_id"`
	Step       int    `json:"step" gorm:"column:step"`
	WfInstID   int    `json:"wf_inst_id" gorm:"column:wf_inst_id"`
	IsFinished bool   `json:"is_finished" gorm:"column:is_finished"`
}

func (w *WfTask) TableName() string {
	return "wf_task"
}
