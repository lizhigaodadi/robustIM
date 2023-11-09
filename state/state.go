package state

import "sync"

type ConnState struct {
	connId         uint64
	AckTimer       *Timer
	ReConnTimer    *Timer
	HeartbeatTimer *Timer /*For Client Use StateServer And BusinessServer Do Not Need*/
	m              sync.Mutex
}

func NewConnState(connId uint64) *ConnState {
	return &ConnState{
		connId:         connId,
		AckTimer:       NewTimer(),
		ReConnTimer:    NewTimer(),
		HeartbeatTimer: NewTimer(),
		m:              sync.Mutex{},
	}
}

func (cs *ConnState) SetAckTimer() {
	/*TODO: Prepare An Ack Timer */
}

func (cs *ConnState) CloseAckTimer() {
	/*TODO: Close Ack Timer*/
}

func (cs *ConnState) SetReConnTimer() {
	/*TODO:Monitor ReConn Event*/
}
func (cs *ConnState) CloseReConnTimer() {
	/*TODO:Monitor ReConn Event*/
}

func (cs *ConnState) SetHeartbeatTimer() {
	/*TODO: Set The HeartbeatTimer*/
}

func (cs *ConnState) CloseHeartbeatTimer() {
	/*TODO: Set The HeartbeatTimer*/

}
