package discov

import (
	"context"
	"fmt"
	"im/common/config"
	"im/common/prpc/discov/etcd"
)

func GetDiscoInstance(ctx context.Context) (Discovery, error) {
	/*判断这个使用了什么作为注册中心*/
	discovName := config.GetPrpcDiscovName()
	switch discovName {
	case "etcd":
		{
			if !etcd.IsEtcdDiscoveryInit() {
				etcd.RunMainRegister(ctx, config.GetEtcdEndPointsForDiscovery())
			}
			return etcd.GetEtcdDiscoveryInstance(), nil
		}
	}
	return nil, fmt.Errorf("获取注册中心实例失败！")
}
