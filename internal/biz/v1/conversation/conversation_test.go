package conversation

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

func TestConversationBiz_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockConversationStore := store.NewMockConversationStore(ctrl)

	conversations := []*model.ConversationM{
		{ID: 1, AgentID: 10, Title: strPtr("会话A")},
		{ID: 2, AgentID: 10, Title: strPtr("会话B")},
	}

	mockStore.EXPECT().Conversation().Return(mockConversationStore)
	mockConversationStore.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(int64(2), conversations, nil)

	b := New(mockStore)
	resp, err := b.List(context.Background(), &apiv1.ListConversationRequest{Offset: 0, Limit: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(2), resp.TotalCount)
	require.Len(t, resp.Conversations, 2)
	assert.Equal(t, int64(1), resp.Conversations[0].Id)
}

func TestConversationBiz_ListMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMessageStore := store.NewMockMessageStore(ctrl)

	messages := []*model.MessageM{
		{ID: 1, ConversationID: 100, Role: model.MessageRoleUser, Content: "你好"},
		{ID: 2, ConversationID: 100, Role: model.MessageRoleAssistant, Content: "你好"},
	}

	mockStore.EXPECT().Message().Return(mockMessageStore)
	mockMessageStore.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(int64(2), messages, nil)

	b := New(mockStore)
	resp, err := b.ListMessages(context.Background(), &apiv1.ListMessageRequest{ConversationId: 100, Offset: 0, Limit: 20})

	require.NoError(t, err)
	assert.Equal(t, int64(2), resp.TotalCount)
	require.Len(t, resp.Messages, 2)
	assert.Equal(t, model.MessageRoleUser, resp.Messages[0].Role)
}

func TestConversationBiz_ListMessages_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockMessageStore := store.NewMockMessageStore(ctrl)

	mockStore.EXPECT().Message().Return(mockMessageStore)
	mockMessageStore.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(int64(0), nil, errors.New("db error"))

	b := New(mockStore)
	_, err := b.ListMessages(context.Background(), &apiv1.ListMessageRequest{ConversationId: 100})

	require.Error(t, err)
}

func TestConversationBiz_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockConversationStore := store.NewMockConversationStore(ctrl)

	mockStore.EXPECT().Conversation().Return(mockConversationStore)
	mockConversationStore.EXPECT().
		Delete(gomock.Any(), where.F("id", int64(1))).
		Return(nil)

	b := New(mockStore)
	_, err := b.Delete(context.Background(), &apiv1.DeleteConversationRequest{Id: 1})

	require.NoError(t, err)
}

func TestConversationBiz_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockIStore(ctrl)
	mockConversationStore := store.NewMockConversationStore(ctrl)

	mockStore.EXPECT().Conversation().Return(mockConversationStore)
	mockConversationStore.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(errors.New("delete failed"))

	b := New(mockStore)
	_, err := b.Delete(context.Background(), &apiv1.DeleteConversationRequest{Id: 1})

	require.Error(t, err)
}
