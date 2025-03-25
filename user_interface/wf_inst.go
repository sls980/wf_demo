package userinterface

import (
	"context"
	"strings"

	"wf_demo/app"
	"wf_demo/infra/util"

	"github.com/spf13/cast"
)

// 发起流程

var (
	wfScheduler app.WfScheduleService
)

func init() {
	wfScheduler = app.NewWfScheduleServiceImpl()
	registerCmd("triggerWf", &CmdWithDesc{
		Handler: triggerWf,
		Desc:    "发起流程",
	})
	registerCmd("completeTask", &CmdWithDesc{
		Handler: completeTask,
		Desc:    "完成任务",
	})
	registerCmd("getAllPendingInst", &CmdWithDesc{
		Handler: getAllPendingInst,
		Desc:    "获取所有待处理的流程",
	})
}

func triggerWf(ctx context.Context, args []string) {
	if len(args) < 2 {
		util.GetLogger(ctx).Infof("usage: triggerWf def_id user_id=<user_id> title=<title> leave_type=<leave_type>")
		return
	}
	wfDefId := cast.ToInt(args[0])
	params := make(map[string]string)
	for _, arg := range args[1:] {
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			util.GetLogger(ctx).Errorf("参数错误")
			return
		}
		params[kv[0]] = kv[1]
	}
	wfInstID, err := wfScheduler.StartProcess(ctx, &app.StartProcessReq{
		WfDefId: wfDefId,
		Params:  params,
	})
	if err != nil {
		println(err.Error())
		return
	}
	util.GetLogger(ctx).Infof("流程发起成功, 流程ID: %d", wfInstID)
}

func completeTask(ctx context.Context, args []string) {
	if len(args) != 2 {
		util.GetLogger(ctx).Infof("usage: completeTask inst_id <true/false>")
		return
	}
	wfInsID := cast.ToInt(args[0])
	pass := cast.ToBool(args[1])
	err := wfScheduler.CompleteTask(ctx, wfInsID, pass)
	if err != nil {
		println(err.Error())
		return
	}
	util.GetLogger(ctx).Infof("任务完成, 流程ID: %d", wfInsID)
}

func getAllPendingInst(ctx context.Context, args []string) {
	wfInsts, err := wfScheduler.GetAllPendingInst(ctx)
	if err != nil {
		println(err.Error())
		return
	}
	for _, item := range wfInsts {
		logger.Printf("title: %s, inst: %d, node_id=%s", item.Title, item.ID, item.NodeID)
	}
}
