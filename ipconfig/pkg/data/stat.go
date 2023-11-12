package data

type Stat struct {
	MessageBytes float64
	ConnectNum   float64
}

func (s *Stat) Add(other *Stat) {
	s.MessageBytes += other.MessageBytes
	s.ConnectNum += other.ConnectNum
}

func (s *Stat) Sub(other *Stat) {
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
