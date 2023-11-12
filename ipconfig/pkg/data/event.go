package data

import (
	"fmt"
	"im/common/discovery"
	"strconv"
)

const (
	AddEvent = 0
	DelEvent = 1
)

/*Represent event type*/
type EventType byte

type Event struct {
	host         string
	port         int
	t            EventType
	connNum      float64
	messageBytes float64
}

func (e *Event) GetHostPort() string {
	return fmt.Sprintf("%s:%d", e.host, e.port)
}

// 设置 host 字段的值
func (e *Event) SetHost(host string) {
	e.host = host
}

// 获取 host 字段的值
func (e *Event) GetHost() string {
	return e.host
}

// 设置 port 字段的值
func (e *Event) SetPort(port int) {
	e.port = port
}

// 获取 port 字段的值
func (e *Event) GetPort() int {
	return e.port
}

// 设置 t 字段的值
func (e *Event) SetEventType(t EventType) {
	e.t = t
}

// 获取 t 字段的值
func (e *Event) GetEventType() EventType {
	return e.t
}

// 设置 connNum 字段的值
func (e *Event) SetConnNum(connNum float64) {
	e.connNum = connNum
}

// 获取 connNum 字段的值
func (e *Event) GetConnNum() float64 {
	return e.connNum
}

// 设置 messageBytes 字段的值
func (e *Event) SetMessageBytes(messageBytes float64) {
	e.messageBytes = messageBytes
}

// 获取 messageBytes 字段的值
func (e *Event) GetMessageBytes() float64 {
	return e.messageBytes
}

func (e *Event) Host() {

}

func NewEvent(info *discovery.EndPointInfo, t EventType) *Event {
	/*TODO:Extract relevant information*/
	host, ok := info.Meta["host"]
	if !ok {
		return nil
	}
	portStr, ok := info.Meta["port"]
	if !ok {
		return nil
	}
	port, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return nil
	}
	messageBytesStr, ok := info.Meta["messageBytes"]
	if !ok {
		return nil
	}
	messageBytes, err := strconv.ParseInt(messageBytesStr, 10, 64)
	if err != nil {
		return nil
	}
	connNumStr, ok := info.Meta["connNum"]
	if !ok {
		return nil
	}
	connNum, err := strconv.ParseInt(connNumStr, 10, 64)
	if err != nil {
		return nil
	}
	if err != nil {
		return nil
	}

	e := &Event{
		host:         host,
		port:         int(port),
		t:            t,
		messageBytes: float64(messageBytes),
		connNum:      float64(connNum),
	}
	return e
}
