package model_provider

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/robinlg/agentops-platform/internal/apiserver/model"
	"github.com/robinlg/agentops-platform/internal/apiserver/store"
	"github.com/robinlg/agentops-platform/internal/pkg/conversion"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/apiserver/v1"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// ModelProviderBiz 定义处理模型提供商请求所需的方法
type ModelProviderBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateModelProviderRequest) (*apiv1.CreateModelProviderResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateModelProviderRequest) (*apiv1.UpdateModelProviderResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteModelProviderRequest) (*apiv1.DeleteModelProviderResponse, error)
	Get(ctx context.Context, rq *apiv1.GetModelProviderRequest) (*apiv1.GetModelProviderResponse, error)
	List(ctx context.Context, rq *apiv1.ListModelProviderRequest) (*apiv1.ListModelProviderResponse, error)

	ModelProviderExpansion
}

// ModelProviderExpansion 定义模型提供商的扩展方法
type ModelProviderExpansion interface {
}

// modelProviderBiz 是 ModelProviderBiz 接口的实现.
type modelProviderBiz struct {
	store store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ ModelProviderBiz = (*modelProviderBiz)(nil)

func New(store store.IStore) *modelProviderBiz {
	return &modelProviderBiz{store: store}
}

// Create 实现 ModelProviderBiz 接口中的 Create 方法
func (b *modelProviderBiz) Create(ctx context.Context, rq *apiv1.CreateModelProviderRequest) (*apiv1.CreateModelProviderResponse, error) {
	var modelProviderM model.ModelProviderM
	_ = copier.Copy(&modelProviderM, rq)

	if err := b.store.ModelProvider().Create(ctx, &modelProviderM); err != nil {
		return nil, err
	}

	return &apiv1.CreateModelProviderResponse{Id: modelProviderM.ID}, nil
}

// Update 实现 ModelProviderBiz 接口中的 Update 方法
func (b *modelProviderBiz) Update(ctx context.Context, rq *apiv1.UpdateModelProviderRequest) (*apiv1.UpdateModelProviderResponse, error) {
	modelProviderM, err := b.store.ModelProvider().Get(ctx, where.F("id", rq.GetId()))
	if err != nil {
		return nil, err
	}

	if err = b.store.ModelProvider().Update(ctx, modelProviderM); err != nil {
		return nil, err
	}

	return &apiv1.UpdateModelProviderResponse{}, nil
}

// Delete 实现 ModelProviderBiz 接口中的 Delete 方法
func (b *modelProviderBiz) Delete(ctx context.Context, rq *apiv1.DeleteModelProviderRequest) (*apiv1.DeleteModelProviderResponse, error) {
	if err := b.store.ModelProvider().Delete(ctx, where.F("id", rq.GetId())); err != nil {
		return nil, err
	}

	return &apiv1.DeleteModelProviderResponse{}, nil
}

// Get 实现 ModelProviderBiz 接口中的 Get 方法
func (b *modelProviderBiz) Get(ctx context.Context, rq *apiv1.GetModelProviderRequest) (*apiv1.GetModelProviderResponse, error) {
	modelProviderM, err := b.store.ModelProvider().Get(ctx, where.F("id", rq.GetId()))
	if err != nil {
		return nil, err
	}

	return &apiv1.GetModelProviderResponse{ModelProvider: conversion.ModelProviderMToModelProviderV1(modelProviderM)}, nil
}

// List 实现 ModelProviderBiz 接口中的 List 方法
func (b *modelProviderBiz) List(ctx context.Context, rq *apiv1.ListModelProviderRequest) (*apiv1.ListModelProviderResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))

	count, modelProviderList, err := b.store.ModelProvider().List(ctx, whr)
	if err != nil {
		return nil, err
	}
	modelProviders := make([]*apiv1.ModelProvider, 0, len(modelProviderList))
	for _, item := range modelProviderList {
		modelProviders = append(modelProviders, conversion.ModelProviderMToModelProviderV1(item))
	}

	return &apiv1.ListModelProviderResponse{TotalCount: count, ModelProviders: modelProviders}, nil
}
