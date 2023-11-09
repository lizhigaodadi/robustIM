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

	Sm = NewStatManager(eventChan)
	/*Enable listening event*/
	go Sm.RunListening()

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
	log.Printf("Monitor Module Init Successfully\n")

}

func sendEvent(val string, t data.EventType) {
	/*Start deserializing*/
	endPointInfo := discovery.UnMarshal([]byte(val))
	if endPointInfo == nil {
		log.Printf("Deserialization failure\n")
		return
	}
	event := data.NewEvent(endPointInfo, t)
	if event == nil {
		log.Printf("New Event Failed: %v\n", event)
		return
	}
	eventChan <- event
}
