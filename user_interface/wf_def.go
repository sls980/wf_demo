package userinterface

import (
	"context"
	"os"

	"wf_demo/app"

	"github.com/spf13/cast"
)

var (
	wfDefSrv app.WfDefService
)

func init() {
	wfDefSrv = app.NewWfDefServiceImpl()
	registerCmd("createWfDef", &CmdWithDesc{
		Handler: createWfDef,
		Desc:    "创建流程定义",
	})
	registerCmd("getWfDef", &CmdWithDesc{
		Handler: getWfDef,
		Desc:    "获取流程定义",
	})
	registerCmd("deleteWfDef", &CmdWithDesc{
		Handler: deleteWfDef,
		Desc:    "删除流程定义",
	})
	registerCmd("getWfDefList", &CmdWithDesc{
		Handler: getWfDefList,
		Desc:    "获取流程定义列表",
	})
}

func createWfDef(ctx context.Context, args []string) {
	if len(args) != 2 {
		logger.Infof("usage: createWfDef name filepath")
		return
	}
	// 从文件中读取流程配置内容
	content, err := os.ReadFile(args[1])
	if err != nil {
		logger.Errorf("read file failed, err: %v", err)
		return
	}
	wfDef := &app.WfDef{
		Name:    args[0],
		Setting: string(content),
	}
	id, err := wfDefSrv.SaveWfDef(ctx, wfDef)
	if err != nil {
		logger.Errorf("save wf def failed, err: %v", err)
		return
	}
	logger.Info(id)
}

func getWfDef(ctx context.Context, args []string) {
	if len(args) != 1 {
		logger.Infof("usage: getWfDef def_id")
		return
	}
	wfDef, err := wfDefSrv.GetWfDefById(ctx, cast.ToInt(args[0]))
	if err != nil {
		logger.Errorf("get wf def failed, err: %v", err)
		return
	}
	logger.Info(wfDef)
}

func getWfDefList(ctx context.Context, args []string) {
	wfDefs, err := wfDefSrv.GetAllWfDef(ctx)
	if err != nil {
		logger.Errorf("get wf def list failed, err: %v", err)
		return
	}
	for _, defItem := range wfDefs {
		logger.Infof("id: %d, name: %s", defItem.Id, defItem.Name)
	}
}

func deleteWfDef(ctx context.Context, args []string) {
	if len(args) != 1 {
		logger.Infof("usage: deleteWfDef def_id")
		return
	}
	err := wfDefSrv.DeleteWfDef(ctx, cast.ToInt(args[0]))
	if err != nil {
		logger.Errorf("delete wf def failed, err: %v", err)
		return
	}
}
