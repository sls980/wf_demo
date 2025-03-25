package userinterface

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"wf_demo/infra/util"

	"github.com/sirupsen/logrus"
)

type CmdFunc func(ctx context.Context, args []string)

type CmdWithDesc struct {
	Handler CmdFunc
	Desc    string
}

var (
	ctx    context.Context
	logger *logrus.Entry
	cmdMap = map[string]*CmdWithDesc{}
)

func init() {
	logger = logrus.NewEntry(logrus.StandardLogger())
	ctx = util.WithLogger(context.Background(), logger)
	registerCmd("help", &CmdWithDesc{
		Handler: printHelp,
		Desc:    "显示帮助信息",
	})
	registerCmd("exit", &CmdWithDesc{
		Handler: func(ctx context.Context, args []string) {
			fmt.Println("退出程序")
			os.Exit(0)
		},
		Desc: "退出程序",
	})
}

func registerCmd(cmd string, fn *CmdWithDesc) {
	cmdMap[cmd] = fn
}

// 运行用户命令行交互
func Run() {
	// 按行处理用户命令行输入
	fmt.Println("欢迎使用工作流演示系统，输入 help 查看帮助，输入 exit 退出")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		// 解析命令行参数
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		cmd := args[0]
		cmdArgs := args[1:]
		// 选择对应的cmd执行
		if fn, ok := cmdMap[cmd]; ok {
			fn.Handler(ctx, cmdArgs)
		} else {
			fmt.Printf("未知命令: %s，输入 help 查看帮助\n", cmd)
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Errorf("读取输入错误: %v", err)
	}
}

// 打印帮助信息
func printHelp(ctx context.Context, args []string) {
	fmt.Println("可用命令:")
	// 找出最长的命令名称长度，用于对齐
	maxCmdLen := 0
	for cmd := range cmdMap {
		if len(cmd) > maxCmdLen {
			maxCmdLen = len(cmd)
		}
	}
	// 按命令名称排序并输出
	var cmds []string
	for cmd := range cmdMap {
		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		if cmd == "help" || cmd == "exit" {
			continue
		}
		fmt.Printf("  %-*s - %s\n", maxCmdLen+2, cmd, cmdMap[cmd].Desc)
	}
	fmt.Printf("  %-*s - %s\n", maxCmdLen+2, "help", cmdMap["help"].Desc)
	fmt.Printf("  %-*s - %s\n", maxCmdLen+2, "exit", cmdMap["exit"].Desc)
}
