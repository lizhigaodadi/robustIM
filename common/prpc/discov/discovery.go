package discov

import (
	"context"
)

type Discovery interface {
	NotifyListeners()
	AddNotify(f func())
	Name() string
	RegisterService(ctx context.Context, service *Service)
	UnRegisterService(ctx context.Context, service *Service)
	GetService(ctx context.Context, serviceName string) *Service
	ShowServiceMessage(ctx context.Context) string /*展示所有注册和发现服务的信息（方便调试）*/
}
