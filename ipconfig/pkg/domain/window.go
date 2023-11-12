package domain

import (
	"im/ipconfig/pkg/data"
)

const (
	DefaultWindowSize = 8
)

type StateWindow struct {
	StatQueue []*data.Stat
	StateChan chan *data.Stat
	idx       int
	queueSize int
	host      string
	port      int
	sumStat   *data.Stat
}

func NewSateWindow(windowSize int, host string, port int, stat *data.Stat) *StateWindow {
	sw := &StateWindow{
		StateChan: make(chan *data.Stat, windowSize),
		idx:       0,
		queueSize: windowSize,
		StatQueue: make([]*data.Stat, windowSize),
		host:      host,
		port:      port,
		sumStat:   stat.Clone(),
	}

	sw.StatQueue = append(sw.StatQueue, stat)
	sw.idx++
	return sw
}

func (sw *StateWindow) NowStat() *data.Stat {
	return sw.StatQueue[sw.idx%sw.queueSize]
}

func (sw *StateWindow) RunListening() error {
	/*TODO:Listen for any incidents in the pipeline*/

	for newStat := range sw.StateChan {
		sw.sumStat.Add(newStat)
		sw.sumStat.Sub(sw.NowStat())
		sw.idx++
	}

	return nil
}

func (sw *StateWindow) PushStat(stat *data.Stat) {
	sw.StateChan <- stat
}

func (sw *StateWindow) GetStaticScore() float64 {
	return sw.sumStat.Avg(sw.queueSize).ConnectNum
}

func (sw *StateWindow) GetDynamicScore() float64 {
	return sw.sumStat.Avg(sw.queueSize).MessageBytes

}
