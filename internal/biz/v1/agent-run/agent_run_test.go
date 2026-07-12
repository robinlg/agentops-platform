package agentrun

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

func TestAgentRunBiz_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentRunStore := store.NewMockAgentRunStore(ctrl)

	agentRuns := []*model.AgentRunM{
		{ID: 1, AgentID: 10, ConversationID: 100, Status: "success"},
		{ID: 2, AgentID: 11, ConversationID: 101, Status: "failed"},
	}

	mockStore.EXPECT().AgentRun().Return(mockAgentRunStore)
	mockAgentRunStore.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(int64(2), agentRuns, nil)

	b := New(mockStore)
	resp, err := b.List(context.Background(), &apiv1.ListAgentRunRequest{Offset: 0, Limit: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(2), resp.TotalCount)
	require.Len(t, resp.AgentRuns, 2)
	assert.Equal(t, int64(1), resp.AgentRuns[0].Id)
	assert.Equal(t, "success", resp.AgentRuns[0].Status)
}

func TestAgentRunBiz_List_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentRunStore := store.NewMockAgentRunStore(ctrl)

	mockStore.EXPECT().AgentRun().Return(mockAgentRunStore)
	mockAgentRunStore.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(int64(0), nil, errors.New("db error"))

	b := New(mockStore)
	_, err := b.List(context.Background(), &apiv1.ListAgentRunRequest{})

	require.Error(t, err)
}

func TestAgentRunBiz_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentRunStore := store.NewMockAgentRunStore(ctrl)

	agentRun := &model.AgentRunM{ID: 5, AgentID: 10, ConversationID: 100, Status: "running"}

	mockStore.EXPECT().AgentRun().Return(mockAgentRunStore)
	mockAgentRunStore.EXPECT().
		Get(gomock.Any(), where.F("id", int64(5))).
		Return(agentRun, nil)

	b := New(mockStore)
	resp, err := b.Get(context.Background(), &apiv1.GetAgentRunRequest{Id: 5})

	require.NoError(t, err)
	require.NotNil(t, resp.AgentRun)
	assert.Equal(t, int64(5), resp.AgentRun.Id)
	assert.Equal(t, "running", resp.AgentRun.Status)
}

func TestAgentRunBiz_Get_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockAgentRunStore := store.NewMockAgentRunStore(ctrl)

	mockStore.EXPECT().AgentRun().Return(mockAgentRunStore)
	mockAgentRunStore.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("not found"))

	b := New(mockStore)
	_, err := b.Get(context.Background(), &apiv1.GetAgentRunRequest{Id: 999})

	require.Error(t, err)
}
