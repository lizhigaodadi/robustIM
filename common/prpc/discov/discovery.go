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
}
