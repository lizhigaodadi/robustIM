package discov

import "im/logger"

type Service struct {
	ServiceName string      `json:"service_name"`
	EndPoints   []*EndPoint `json:"end_points"`
}

type EndPoint struct {
	ServerName string `json:"server_name"`
	Ip         string `json:"ip"`
	Port       uint32 `json:"port"`
	Weight     uint32 `json:"weight"`
}

func (s *Service) AddService(other *Service) {
	/*确保这两个服务是一样的*/
	if s.ServiceName != other.ServiceName {
		logger.StdLog().Warnf("Add A Different To The %s Service", s.ServiceName)
		return
	}

	for _, ep := range other.EndPoints {
		var isAdd bool = true
		for _, e := range s.EndPoints {
			/*判断一下是否相等*/
			if ep.Equals(e) {
				isAdd = false
				break
			}
		}

		if isAdd {
			s.EndPoints = append(s.EndPoints, ep)
		}
	}

}

func (s *Service) RemoveService(other *Service) {

	if s.ServiceName != other.ServiceName {
		logger.StdLog().Warnf("Remove A Different To The %s Service", s.ServiceName)
		return
	}

	n := make([]*EndPoint, 0, len(s.EndPoints))

	for _, ep := range s.EndPoints {
		var isRemove bool = false
		for _, e := range other.EndPoints {
			if ep.Equals(e) {
				isRemove = true
				break
			}
		}
		if !isRemove {
			n = append(n, ep)
		}
	}
	s.EndPoints = n

}

func (e *EndPoint) Equals(o *EndPoint) bool {
	if e.Ip == o.Ip && e.Port == o.Port && e.ServerName == o.ServerName {
		return true
	}
	return false
}

func (e *EndPoint) EqualsAndUpdate(o *EndPoint) bool {
	if e.Equals(o) {
		e.Weight = o.Weight
		return true
	} else {
		return false
	}
}
