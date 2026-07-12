package agentrun

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/pkg/conversion"
	"github.com/robinlg/agentops-platform/internal/store"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// AgentRunBiz 定义处理智能体运行记录请求所需的方法
type AgentRunBiz interface {
	List(ctx context.Context, rq *apiv1.ListAgentRunRequest) (*apiv1.ListAgentRunResponse, error)
	Get(ctx context.Context, rq *apiv1.GetAgentRunRequest) (*apiv1.GetAgentRunResponse, error)

	AgentRunExpansion
}

// AgentRunExpansion 定义智能体运行记录的扩展方法
type AgentRunExpansion interface {
}

// agentRunBiz 是 AgentRunBiz 接口的实现.
type agentRunBiz struct {
	store store.IStore
}

// 确保 agentRunBiz 实现了 AgentRunBiz 接口
var _ AgentRunBiz = (*agentRunBiz)(nil)

func New(store store.IStore) *agentRunBiz {
	return &agentRunBiz{store: store}
}

// List 实现 AgentRunBiz 接口中的 List 方法
func (b *agentRunBiz) List(ctx context.Context, rq *apiv1.ListAgentRunRequest) (*apiv1.ListAgentRunResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))

	count, agentRunList, err := b.store.AgentRun().List(ctx, whr)
	if err != nil {
		return nil, err
	}
	agentRuns := make([]*apiv1.AgentRun, 0, len(agentRunList))
	for _, item := range agentRunList {
		agentRuns = append(agentRuns, conversion.AgentRunMToAgentRunV1(item))
	}

	return &apiv1.ListAgentRunResponse{TotalCount: count, AgentRuns: agentRuns}, nil
}

// Get 实现 AgentRunBiz 接口中的 Get 方法
func (b *agentRunBiz) Get(ctx context.Context, rq *apiv1.GetAgentRunRequest) (*apiv1.GetAgentRunResponse, error) {
	agentRunM, err := b.store.AgentRun().Get(ctx, where.F("id", rq.GetId()))
	if err != nil {
		return nil, err
	}

	return &apiv1.GetAgentRunResponse{AgentRun: conversion.AgentRunMToAgentRunV1(agentRunM)}, nil
}
