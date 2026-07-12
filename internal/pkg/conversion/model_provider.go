package conversion

import (
	"github.com/robinlg/agentops-platform/internal/apiserver/model"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/apiserver/v1"
	"github.com/robinlg/onexlib/pkg/core"
)

// ModelProviderMToModelProviderV1 将模型层的 ModelProviderM（模型提供商模型对象）转换为 Protobuf 层的 ModelProvider（v1 模型提供商对象）
func ModelProviderMToModelProviderV1(modelProviderModel *model.ModelProviderM) *apiv1.ModelProvider {
	var protoModelProvider apiv1.ModelProvider
	_ = core.CopyWithConverters(&protoModelProvider, modelProviderModel)
	return &protoModelProvider
}
