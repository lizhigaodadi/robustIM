package state

import (
	"im/common/cache"
	"sync"
)

var cacheState *CacheState

type CacheState struct {
	msgId            uint32
	ConnToStateTable sync.Map
}

func CacheInit() {
	cacheState = &CacheState{
		msgId:            0,
		ConnToStateTable: sync.Map{},
	}
	cache.Init() /*Init Redis*/
}

func (cs *CacheState) AddConn(state *ConnState) {
	cs.ConnToStateTable.Store(state.connId, state)
}

/*dId + userId*/
func (cs *CacheState) AddLoginState(connId uint64) {

}

func (cs *CacheState) GetConnState(connId uint64) *ConnState {
	connState, ok := cs.ConnToStateTable.Load(connId)
	if !ok {
		return nil
	}

	state, ok := connState.(*ConnState)
	if !ok {
		panic("Not Match Type You Want To Get From CacheState")
	}

	return state
}
