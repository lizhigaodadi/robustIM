package domain

import (
	"context"
	"im/common/config"
	"im/common/discovery"
	"im/ipconfig/pkg/data"
	"log"
)

var (
	eventChan     chan *data.Event
	eventChanSize = 5
)

/*Init Monitor To Etcd*/
func Init(ctx *context.Context) {
	eventChan = make(chan *data.Event, eventChanSize)
	/*Init Etcd Discovery Service*/

	Ms = NewStatManager(eventChan)
	/*Enable listening event*/
	go Ms.RunListening()

	dis := discovery.NewEtcdServiceDiscovery(ctx)
	addNode := func(key, val string) {
		sendEvent(val, data.AddEvent)
	}
	delNode := func(key, val string) {
		sendEvent(val, data.DelEvent)
	}

	err := dis.WatchService(config.GetIpConfigPath(), addNode, delNode)
	if err != nil {
		log.Fatalf("Init Watch Etcd Service Failed\n")
	}

}

func sendEvent(val string, t data.EventType) {
	/*Start deserializing*/
	endPointInfo := discovery.UnMarshal([]byte(val))
	event := data.NewEvent(endPointInfo, data.EventType(t))
	eventChan <- event
}
