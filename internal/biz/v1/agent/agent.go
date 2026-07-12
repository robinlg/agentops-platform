package agent

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/robinlg/agentops-platform/internal/model"
	"github.com/robinlg/agentops-platform/internal/pkg/conversion"
	"github.com/robinlg/agentops-platform/internal/store"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// AgentBiz 定义处理智能体请求所需的方法
type AgentBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateAgentRequest) (*apiv1.CreateAgentResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateAgentRequest) (*apiv1.UpdateAgentResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteAgentRequest) (*apiv1.DeleteAgentResponse, error)
	Get(ctx context.Context, rq *apiv1.GetAgentRequest) (*apiv1.GetAgentResponse, error)
	List(ctx context.Context, rq *apiv1.ListAgentRequest) (*apiv1.ListAgentResponse, error)

	AgentExpansion
}

// AgentExpansion 定义智能体的扩展方法
type AgentExpansion interface {
}

// agentBiz 是 AgentBiz 接口的实现.
type agentBiz struct {
	store store.IStore
}

// 确保 agentBiz 实现了 AgentBiz 接口
var _ AgentBiz = (*agentBiz)(nil)

func New(store store.IStore) *agentBiz {
	return &agentBiz{store: store}
}

// Create 实现 AgentBiz 接口中的 Create 方法
func (b *agentBiz) Create(ctx context.Context, rq *apiv1.CreateAgentRequest) (*apiv1.CreateAgentResponse, error) {
	var agentM model.AgentM
	_ = copier.Copy(&agentM, rq)

	if err := b.store.Agent().Create(ctx, &agentM); err != nil {
		return nil, err
	}

	return &apiv1.CreateAgentResponse{Id: agentM.ID}, nil
}

// Update 实现 AgentBiz 接口中的 Update 方法
func (b *agentBiz) Update(ctx context.Context, rq *apiv1.UpdateAgentRequest) (*apiv1.UpdateAgentResponse, error) {
	agentM, err := b.store.Agent().Get(ctx, where.F("id", rq.GetId()))
	if err != nil {
		return nil, err
	}

	if rq.Name != nil {
		agentM.Name = *rq.Name
	}
	if rq.Description != nil {
		agentM.Description = rq.Description
	}
	if rq.SystemPrompt != nil {
		agentM.SystemPrompt = rq.SystemPrompt
	}
	if rq.ModelProviderId != nil {
		agentM.ModelProviderID = *rq.ModelProviderId
	}
	if rq.Model != nil {
		agentM.Model = rq.Model
	}
	if rq.Temperature != nil {
		agentM.Temperature = rq.Temperature
	}
	if rq.MaxTokens != nil {
		agentM.MaxTokens = rq.MaxTokens
	}

	if err = b.store.Agent().Update(ctx, agentM); err != nil {
		return nil, err
	}

	return &apiv1.UpdateAgentResponse{}, nil
}

// Delete 实现 AgentBiz 接口中的 Delete 方法
func (b *agentBiz) Delete(ctx context.Context, rq *apiv1.DeleteAgentRequest) (*apiv1.DeleteAgentResponse, error) {
	if err := b.store.Agent().Delete(ctx, where.F("id", rq.GetId())); err != nil {
		return nil, err
	}

	return &apiv1.DeleteAgentResponse{}, nil
}

// Get 实现 AgentBiz 接口中的 Get 方法
func (b *agentBiz) Get(ctx context.Context, rq *apiv1.GetAgentRequest) (*apiv1.GetAgentResponse, error) {
	agentM, err := b.store.Agent().Get(ctx, where.F("id", rq.GetId()))
	if err != nil {
		return nil, err
	}

	return &apiv1.GetAgentResponse{Agent: conversion.AgentMToAgentV1(agentM)}, nil
}

// List 实现 AgentBiz 接口中的 List 方法
func (b *agentBiz) List(ctx context.Context, rq *apiv1.ListAgentRequest) (*apiv1.ListAgentResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))

	count, agentList, err := b.store.Agent().List(ctx, whr)
	if err != nil {
		return nil, err
	}
	agents := make([]*apiv1.Agent, 0, len(agentList))
	for _, item := range agentList {
		agents = append(agents, conversion.AgentMToAgentV1(item))
	}

	return &apiv1.ListAgentResponse{TotalCount: count, Agents: agents}, nil
}
