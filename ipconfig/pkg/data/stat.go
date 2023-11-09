package data

type Stat struct {
	MessageBytes float64
	ConnectNum   float64
}

func (s *Stat) Add(other *Stat) {
	if other == nil {
		return
	}

	s.MessageBytes += other.MessageBytes
	s.ConnectNum += other.ConnectNum
}

func (s *Stat) Sub(other *Stat) {
	if other == nil {
		return
	}
	s.MessageBytes -= other.MessageBytes
	s.ConnectNum -= other.ConnectNum
}

func (s *Stat) Clone() *Stat {
	return &Stat{
		MessageBytes: s.MessageBytes,
		ConnectNum:   s.ConnectNum,
	}
}

func (s *Stat) Avg(num int) *Stat {
	other := s.Clone()
	s.MessageBytes /= float64(num)
	s.ConnectNum /= float64(num)
	return other
}

func (s *Stat) IsEnd() bool {
	if s.MessageBytes == -44.44 {
		return true
	}
	return false
}

func NewEndStat() *Stat {
	return &Stat{
		MessageBytes: 44.44,
	}
}
