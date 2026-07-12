package model_provider

import (
	"context"
	"errors"
	"testing"

	"github.com/robinlg/agentops-platform/internal/model"
	"github.com/robinlg/agentops-platform/internal/store"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/robinlg/onexlib/pkg/store/where"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func strPtr(s string) *string { return &s }

func TestModelProviderBiz_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMPStore := store.NewMockModelProviderStore(ctrl)

	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, mp *model.ModelProviderM) error {
			mp.ID = 88
			return nil
		})

	b := New(mockStore)
	resp, err := b.Create(context.Background(), &apiv1.CreateModelProviderRequest{
		Name:         "openai",
		ProviderType: "openai",
		BaseUrl:      "https://api.openai.com/v1",
		ApiKey:       "sk-xxx",
		DefaultModel: "gpt-4o-mini",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(88), resp.Id)
}

func TestModelProviderBiz_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMPStore := store.NewMockModelProviderStore(ctrl)

	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	b := New(mockStore)
	_, err := b.Create(context.Background(), &apiv1.CreateModelProviderRequest{Name: "x"})

	require.Error(t, err)
}

func TestModelProviderBiz_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMPStore := store.NewMockModelProviderStore(ctrl)

	existing := &model.ModelProviderM{ID: 1, Name: "old", DefaultModel: "gpt-3.5"}

	// Get 返回已有记录
	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().Get(gomock.Any(), where.F("id", int64(1))).Return(existing, nil)

	// Update 校验字段已被部分更新
	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, mp *model.ModelProviderM) error {
			assert.Equal(t, "new-name", mp.Name)
			assert.Equal(t, "gpt-4o", mp.DefaultModel)
			return nil
		})

	b := New(mockStore)
	_, err := b.Update(context.Background(), &apiv1.UpdateModelProviderRequest{
		Id:           1,
		Name:         strPtr("new-name"),
		DefaultModel: strPtr("gpt-4o"),
	})

	require.NoError(t, err)
}

func TestModelProviderBiz_Update_GetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMPStore := store.NewMockModelProviderStore(ctrl)

	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found"))

	b := New(mockStore)
	_, err := b.Update(context.Background(), &apiv1.UpdateModelProviderRequest{Id: 999})

	require.Error(t, err)
}

func TestModelProviderBiz_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMPStore := store.NewMockModelProviderStore(ctrl)

	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().Delete(gomock.Any(), where.F("id", int64(1))).Return(nil)

	b := New(mockStore)
	_, err := b.Delete(context.Background(), &apiv1.DeleteModelProviderRequest{Id: 1})

	require.NoError(t, err)
}

func TestModelProviderBiz_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMPStore := store.NewMockModelProviderStore(ctrl)

	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().Get(gomock.Any(), where.F("id", int64(1))).
		Return(&model.ModelProviderM{ID: 1, Name: "openai"}, nil)

	b := New(mockStore)
	resp, err := b.Get(context.Background(), &apiv1.GetModelProviderRequest{Id: 1})

	require.NoError(t, err)
	require.NotNil(t, resp.ModelProvider)
	assert.Equal(t, int64(1), resp.ModelProvider.Id)
	assert.Equal(t, "openai", resp.ModelProvider.Name)
}

func TestModelProviderBiz_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMPStore := store.NewMockModelProviderStore(ctrl)

	list := []*model.ModelProviderM{
		{ID: 1, Name: "openai"},
		{ID: 2, Name: "deepseek"},
	}

	mockStore.EXPECT().ModelProvider().Return(mockMPStore)
	mockMPStore.EXPECT().List(gomock.Any(), gomock.Any()).Return(int64(2), list, nil)

	b := New(mockStore)
	resp, err := b.List(context.Background(), &apiv1.ListModelProviderRequest{Offset: 0, Limit: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(2), resp.TotalCount)
	require.Len(t, resp.ModelProviders, 2)
	assert.Equal(t, "deepseek", resp.ModelProviders[1].Name)
}
