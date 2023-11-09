package gateWay

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var Id uint32

var generator *ConnIDGenerator

var (
	version      = uint64(0)
	sequenceBits = uint64(16)
	maxSequence  = int64(-1) ^ (int64(-1) << sequenceBits)
	twepoch      = int64(1589923200000)
	timeLeft     = uint8(16)
	versionLeft  = uint8(63)
)

type connection struct {
	Id   uint64 /*Device Id*/
	fd   int
	e    *epoll
	conn *net.TCPConn
}

func (c *connection) Close() {
	ep.id2Conn.Delete(c.Id)
	if c.e != nil {
		c.e.fd2Conn.Delete(c.fd)
	}
	err := c.Close
	panic(err)
}

func (c *connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func generatorInit() {
	generator = &ConnIDGenerator{
		mu:        sync.Mutex{},
		LastStamp: 0,
		Sequence:  0,
	}
}

/*Set TCP connection Keep-alive*/

func SetTCPConfig(conn *net.TCPConn) {
	_ = conn.SetKeepAlive(true)
}

func (c *connection) BindEpoll(e *epoll) {
	c.e = e
}

func NewConnection(conn *net.TCPConn) *connection {
	/**/
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	fd := int(pfdVal.FieldByName("Sysfd").Int())
	did, err := generator.NextId()
	if err != nil {
		panic(err)
	}

	return &connection{
		Id:   did,
		fd:   fd,
		conn: conn,
	}
}

func GetDid() uint32 {
	return atomic.AddUint32(&Id, 1)
}

type ConnIDGenerator struct {
	mu        sync.Mutex
	LastStamp int64 /*Record The Last Time Stamp*/
	Sequence  int64 /*Current ms Has Generated ID Sequence 1ms */
}

func (c *ConnIDGenerator) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (c *ConnIDGenerator) NextId() (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.nextId()
}

func (c *ConnIDGenerator) nextId() (uint64, error) {
	timeStamp := c.getMilliSeconds()
	if timeStamp < c.LastStamp {
		return 0, errors.New("time is moving backwards, waiting until")
	}

	if c.LastStamp == timeStamp {
		c.Sequence = (c.Sequence + 1) & maxSequence
		if c.Sequence == 0 {
			for timeStamp <= c.LastStamp {
				timeStamp = c.getMilliSeconds()
			}
		}
	} else {
		c.Sequence = 0
	}
	c.LastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft)
	connId := uint64(id) | (version << versionLeft)

	return connId, nil
}
