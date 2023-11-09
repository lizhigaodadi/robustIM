package gateWay

import (
	"fmt"
	"golang.org/x/sys/unix"
	"im/common/config"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
)

var TCPNum int32
var ep *EPool

type connectionHandler func(connect *connection)

type EPool struct {
	ln         *net.TCPListener
	doneChan   chan struct{}
	connChan   chan *connection
	epollCount int
	id2Conn    sync.Map
	handler    connectionHandler
}

type epoll struct {
	epollFd int /*epoll 's fd*/
	fd2Conn sync.Map
}

func (ep *epoll) wait(msec int) []*connection {
	events := make([]unix.EpollEvent, config.GetGateWayQueueSize())
	n, err := unix.EpollWait(ep.epollFd, events, msec)
	if err != nil {
		log.Printf("Epoll Wait too long\n")
		return make([]*connection, 0)
	}
	log.Printf("n: %d\n", n)
	res := make([]*connection, n)
	for i := 0; i < n; i++ {
		fd := events[i].Fd
		log.Printf("Load : %d", fd)
		c, ok := ep.fd2Conn.Load(fd)
		if !ok {
			log.Printf("Load Fd To Connection Failed\n")
			continue
		}
		conn, ok := c.(*connection)
		if !ok {
			log.Printf("Convert To Connection Failed\n")
			continue
		}
		res = append(res, conn)
	}

	return res
}

func NewEPool(ln *net.TCPListener, epollCount int, handler connectionHandler) *EPool {

	return &EPool{
		epollCount: epollCount,
		ln:         ln,
		doneChan:   make(chan struct{}, epollCount),
		connChan:   make(chan *connection, epollCount),
		id2Conn:    sync.Map{},
		handler:    handler,
	}
}

func RunEpoll(ln *net.TCPListener) {
	ep = NewEPool(ln, config.GetGateWayEpollCount(), epollHandler)
	ep.ProcessAll()
}

func epollHandler(connect *connection) {
	/*TODO: epoll handler logic*/
	if connect == nil {
		log.Printf("connect is nil\n")
		return
	}
	log.Printf("Has Some Event Happend! fd: %d\n", connect.fd)

}

func NewEpoll() *epoll {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		log.Printf("Create Epoll Failed\n")
		return nil
	}

	return &epoll{
		epollFd: fd,
		fd2Conn: sync.Map{},
	}
}

func (ep *epoll) Add(connect *connection) error {
	log.Printf("Store Connection Fd: %d", connect.fd)
	ep.fd2Conn.Store(int32(connect.fd), connect)
	_, ok := ep.fd2Conn.Load(int32(connect.fd))
	if !ok {
		log.Printf("gggggggggggggggggggggggggggggggggg\n")
	}

	err := unix.EpollCtl(ep.epollFd, unix.EPOLL_CTL_ADD, connect.fd, &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLET,
		Fd:     int32(connect.fd),
	})
	if err != nil {
		log.Printf("Epoll Add event Failed\n")
		return err
	}

	/* Update Message*/

	return nil
}

func (e *EPool) ProcessAll() {
	e.createAcceptProcess()
	epCount := config.GetGateWayEpollCount()
	for i := 0; i < epCount; i++ {
		go e.ProcEpoll()
	}
}

func (e *EPool) ProcEpoll() {

	ep := NewEpoll()

	go func() {
		for {
			select {
			case <-e.doneChan:
				{
					break /*End*/
				}
			case conn := <-e.connChan:
				{
					log.Printf("Add Connections Happend!\n")
					/*Add*/
					err := ep.Add(conn)
					if err != nil {
						log.Printf("Add Connections To EPoll Failed\n")
					}
				}
			}
		}
	}()

	for {
		conns := ep.wait(200)
		for _, conn := range conns {
			wp.Submit(func() { /*Asynchronous operation*/
				e.handler(conn)
			})
		}
	}
}

func (e *EPool) createAcceptProcess() {

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, err := e.ln.AcceptTCP()
				/*Check Our TCP Num*/
				if !checkTCP() {
					/*Over Connections Limit*/
					_ = conn.Close() /*Connection Close*/
					log.Println("GateWay Has Touch The Connections Limit")
					continue
				}
				if err != nil {
					if ne, ok := err.(net.Error); ok && ne.Temporary() {
						fmt.Errorf("accept temp err: %v\n", ne)
						continue
					}
					fmt.Errorf("accpet err:%v\n", err)
				}
				log.Printf("An New Connection Enter!\n")

				connect := NewConnection(conn)
				log.Printf("Store conn Id: %d\n", connect.Id)
				e.id2Conn.Store(connect.Id, connect)
				e.PushTask(connect)
			}
		}()
	}
}

func (e *EPool) PushTask(conn *connection) {
	e.connChan <- conn
}

/*-------------control TCPNum------------*/

func IncrTCPNum() {
	atomic.AddInt32(&TCPNum, 1)
}
func DecrTCPNum() {
	atomic.AddInt32(&TCPNum, -1)
}
func GetTCPNum() int32 {
	tcpNum := atomic.LoadInt32(&TCPNum)
	return tcpNum
}

func checkTCP() bool {
	num := GetTCPNum()
	bearLimit := config.GetGateConnLimit()

	return num <= bearLimit
}
