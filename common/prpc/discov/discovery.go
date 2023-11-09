package discov

import (
	"context"
)

type Discovery interface {
	NotifyListeners()
	AddNotify(f func())
	Name() string
	RegisterService(ctx context.Context, service *service)
	UnRegisterService(ctx context.Context, service *service)
}
