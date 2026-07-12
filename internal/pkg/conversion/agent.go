package conversion

import (
	"github.com/robinlg/agentops-platform/internal/apiserver/model"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/apiserver/v1"
	"github.com/robinlg/onexlib/pkg/core"
)

// AgentMToAgentV1 将模型层的 AgentM（智能体模型对象）转换为 API 层的 Agent（v1 智能体对象）
func AgentMToAgentV1(agentModel *model.AgentM) *apiv1.Agent {
	var protoAgent apiv1.Agent
	_ = core.CopyWithConverters(&protoAgent, agentModel)
	return &protoAgent
}
