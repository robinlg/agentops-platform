package conversion

import (
	"github.com/robinlg/agentops-platform/internal/model"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/robinlg/onexlib/pkg/core"
)

// AgentRunMToAgentRunV1 将模型层的 AgentRunM 转换为 API 层的 AgentRun
func AgentRunMToAgentRunV1(agentRunModel *model.AgentRunM) *apiv1.AgentRun {
	var protoAgentRun apiv1.AgentRun
	_ = core.CopyWithConverters(&protoAgentRun, agentRunModel)
	return &protoAgentRun
}
