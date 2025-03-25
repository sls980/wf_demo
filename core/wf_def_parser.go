package core

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"wf_demo/infra/util"
)

// 数据结构定义

type NodeType int

const (
	START NodeType = iota + 1
	ROUTE
	CONDITION
	APPROVER
	NOTIFIER
)

type CondType string

const (
	IsAnyOf CondType = "isAnyOf"
	Range   CondType = "range"
	Equal   CondType = "equal"
)

// 流程中的一个节点
type Node struct {
	Name           string       `json:"name"`
	Type           NodeType     `json:"type"`
	NodeID         string       `json:"node_id"`
	PrevID         string       `json:"prev_id"`
	ChildNode      *Node        `json:"child_node"`
	ConditionNodes []*Node      `json:"condition_nodes"`
	Setting        *NodeSetting `json:"setting"`
}

type NodeSetting struct {
	Conditions  []*NodeCondition `json:"conditions"`
	ActionRules []*ActionRule    `json:"action_rules"`
}

type NodeCondition struct {
	Type     CondType `json:"type"`
	ParamKey string   `json:"param_key"`
	// for equal
	Value string `json:"value"`
	// for isAnyOf
	Values []string `json:"values"`
	// for range
	Min int `json:"min"`
	Max int `json:"max"`
}

type ActionRule struct {
	Type      string `json:"type"`
	LabelName string `json:"label_name"`
}

// 节点信息
type NodeInfo struct {
	NodeID       string   `json:"node_id"`
	Type         string   `json:"type"`
	Approver     string   `json:"approver"`
	ApproverType NodeType `json:"approver_type"`
}

func (n *Node) add2ExecutionList(ctx context.Context, list *list.List) {
	switch n.Type {
	case APPROVER, NOTIFIER:
		list.PushBack(NodeInfo{
			NodeID:       n.NodeID,
			Type:         n.Setting.ActionRules[0].Type,
			Approver:     n.Setting.ActionRules[0].LabelName,
			ApproverType: n.Type,
		})
		util.GetLogger(ctx).Infof("add node to execution list, node: %v", n)
	default:
	}
}

type WfDefParser interface {
	ParseWfDef(ctx context.Context, wfSetting string, variable map[string]string) ([]*NodeInfo, error)
}

type WfDefParserImpl struct{}

func NewWfDefParser() WfDefParser {
	return &WfDefParserImpl{}
}

// 流程定义解析器
func (wdp *WfDefParserImpl) ParseWfDef(ctx context.Context, wfSetting string, variable map[string]string) ([]*NodeInfo, error) {
	var rootNode *Node
	err := json.Unmarshal([]byte(wfSetting), &rootNode)
	if err != nil {
		util.GetLogger(ctx).Errorf("parse wf def failed, err: %v", err)
		return nil, err
	}
	nodeList := list.New()
	err = parseNode(ctx, rootNode, variable, nodeList)
	if err != nil {
		util.GetLogger(ctx).Errorf("parse node failed, err: %v", err)
		return nil, err
	}
	// 插入开始&结束节点
	nodeList.PushBack(NodeInfo{
		NodeID: "end",
	})
	nodeList.PushFront(NodeInfo{
		NodeID: "start",
	})
	// 将list转换为数组
	var nodeInfoList []*NodeInfo
	err = util.Decode(util.List2Array(nodeList), &nodeInfoList)
	if err != nil {
		util.GetLogger(ctx).Errorf("json unmarshal失败: %v", err)
		return nil, err
	}
	return nodeInfoList, nil
}

func parseNode(ctx context.Context, node *Node, variable map[string]string, nodeList *list.List) error {
	node.add2ExecutionList(ctx, nodeList)
	// 先处理条件节点
	if len(node.ConditionNodes) > 0 {
		condNode, err := getConditionNode(ctx, node.ConditionNodes, variable)
		if err != nil {
			util.GetLogger(ctx).Errorf("get condition node failed, err: %v", err)
			return err
		}
		if condNode == nil {
			// 如果没有符合条件的节点，返回错误
			return fmt.Errorf("no condition node found %s", node.Name)
		}
		// 递归处理条件节点
		err = parseNode(ctx, condNode, variable, nodeList)
		if err != nil {
			util.GetLogger(ctx).Errorf("parse condiditon node failed, err: %v", err)
			return err
		}
	}
	// 再处理子节点
	if node.ChildNode != nil {
		// 递归处理子节点
		err := parseNode(ctx, node.ChildNode, variable, nodeList)
		if err != nil {
			util.GetLogger(ctx).Errorf("parse child node failed, err: %v", err)
			return err
		}
	}
	return nil
}

func getConditionNode(ctx context.Context, nodes []*Node, variable map[string]string) (*Node, error) {
	for _, node := range nodes {
		// 判断条件是否全部成立
		passCnt := 0
		for _, cond := range node.Setting.Conditions {
			paramValue := variable[cond.ParamKey]
			pass, err := evalCondition(ctx, cond, paramValue)
			if err != nil {
				util.GetLogger(ctx).Errorf("eval condition failed, err: %v", err)
				return nil, err
			}
			if pass {
				passCnt++
			}
		}
		if passCnt == len(node.Setting.Conditions) {
			return node, nil
		}
	}
	return nil, nil
}

// 简单的规则引擎
func evalCondition(ctx context.Context, cond *NodeCondition, value string) (bool, error) {
	// 根据条件类型进行判断，目前只支持equal、range和isAnyOf两种运算符
	switch cond.Type {
	case Equal:
		return value == cond.Value, nil
	case IsAnyOf:
		for _, v := range cond.Values {
			if v == value {
				return true, nil
			}
		}
		return false, nil
	case Range:
		num, err := strconv.Atoi(value)
		if err != nil {
			util.GetLogger(ctx).Errorf("convert value to int failed, err: %v", err)
			return false, err
		}
		return num >= cond.Min && num <= cond.Max, nil
	default:
		return false, nil
	}
}
