package discov

import "im/logger"

type Service struct {
	serviceName string
	endPoints   []*EndPoint
}

type EndPoint struct {
	serverName string
	Ip         string
	Port       uint32
	weight     uint32
}

func (s *Service) AddService(other *Service) {
	/*确保这两个服务是一样的*/
	if s.serviceName != other.serviceName {
		logger.StdLog().Warnf("Add A Different To The %s Service", s.serviceName)
		return
	}

	for _, ep := range s.endPoints {
		var isAdd bool = false
		for _, e := range other.endPoints {
			/*判断一下是否相等*/
		}
	}

}

func (e *EndPoint) Equals(o *EndPoint) bool {
	if e.Ip == o.Ip && e.Port == o.Port && e.serverName == o.serverName {
		return true
	}
	return false
}

func (e *EndPoint) EqualsAndUpdate(o *EndPoint) bool {
	if e.Equals(o) {
		e.weight = o.weight
		return true
	} else {
		return false
	}
}
