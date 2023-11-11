package discov

import (
	"fmt"
	"im/logger"
)

var (
	prpcPrefix = "/prpc/service"
)

func WithPrpcPrefix(serverName string, serviceName string) string {

	return prpcPrefix + "/" + serviceName + "/" + serviceName
}

func NewEndPoints(eps ...*EndPoint) []*EndPoint {
	endPoints := make([]*EndPoint, 0, len(eps))
	for _, ep := range eps {
		endPoints = append(endPoints, ep)
	}

	return endPoints
}

/*将抽象的EndPoint对象转化为Etcd所需要的格式*/

func EndPointsToStrings(eps []*EndPoint) []string {
	strings := make([]string, 0, len(eps))

	for _, e := range eps {
		strings = append(strings, EndPointToString(e))
	}
	return strings
}

func EndPointToString(ep *EndPoint) string {
	if ep == nil {
		logger.Debugf("EndPoint 判断为nil，请检查代码")
		return ""
	}

	return fmt.Sprintf("%s:%d", ep.Ip, ep.Port)
}
