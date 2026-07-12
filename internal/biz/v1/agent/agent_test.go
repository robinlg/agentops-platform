package agent

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
func int64Ptr(v int64) *int64 { return &v }

func TestAgentBiz_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentStore := store.NewMockAgentStore(ctrl)

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, a *model.AgentM) error {
			a.ID = 66
			return nil
		})

	b := New(mockStore)
	resp, err := b.Create(context.Background(), &apiv1.CreateAgentRequest{
		Name:            "assistant",
		ModelProviderId: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(66), resp.Id)
}

func TestAgentBiz_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentStore := store.NewMockAgentStore(ctrl)

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	b := New(mockStore)
	_, err := b.Create(context.Background(), &apiv1.CreateAgentRequest{Name: "x"})

	require.Error(t, err)
}

func TestAgentBiz_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentStore := store.NewMockAgentStore(ctrl)

	existing := &model.AgentM{ID: 1, Name: "old", ModelProviderID: 1}

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().Get(gomock.Any(), where.F("id", int64(1))).Return(existing, nil)

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, a *model.AgentM) error {
			assert.Equal(t, "new-name", a.Name)
			require.NotNil(t, a.Model)
			assert.Equal(t, "gpt-4o", *a.Model)
			assert.Equal(t, int64(2), a.ModelProviderID)
			return nil
		})

	b := New(mockStore)
	_, err := b.Update(context.Background(), &apiv1.UpdateAgentRequest{
		Id:              1,
		Name:            strPtr("new-name"),
		Model:           strPtr("gpt-4o"),
		ModelProviderId: int64Ptr(2),
	})

	require.NoError(t, err)
}

func TestAgentBiz_Update_GetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentStore := store.NewMockAgentStore(ctrl)

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found"))

	b := New(mockStore)
	_, err := b.Update(context.Background(), &apiv1.UpdateAgentRequest{Id: 999})

	require.Error(t, err)
}

func TestAgentBiz_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentStore := store.NewMockAgentStore(ctrl)

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().Delete(gomock.Any(), where.F("id", int64(1))).Return(nil)

	b := New(mockStore)
	_, err := b.Delete(context.Background(), &apiv1.DeleteAgentRequest{Id: 1})

	require.NoError(t, err)
}

func TestAgentBiz_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentStore := store.NewMockAgentStore(ctrl)

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().Get(gomock.Any(), where.F("id", int64(1))).
		Return(&model.AgentM{ID: 1, Name: "assistant"}, nil)

	b := New(mockStore)
	resp, err := b.Get(context.Background(), &apiv1.GetAgentRequest{Id: 1})

	require.NoError(t, err)
	require.NotNil(t, resp.Agent)
	assert.Equal(t, int64(1), resp.Agent.Id)
	assert.Equal(t, "assistant", resp.Agent.Name)
}

func TestAgentBiz_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentStore := store.NewMockAgentStore(ctrl)

	list := []*model.AgentM{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "b"},
	}

	mockStore.EXPECT().Agent().Return(mockAgentStore)
	mockAgentStore.EXPECT().List(gomock.Any(), gomock.Any()).Return(int64(2), list, nil)

	b := New(mockStore)
	resp, err := b.List(context.Background(), &apiv1.ListAgentRequest{Offset: 0, Limit: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(2), resp.TotalCount)
	require.Len(t, resp.Agents, 2)
	assert.Equal(t, "b", resp.Agents[1].Name)
}
