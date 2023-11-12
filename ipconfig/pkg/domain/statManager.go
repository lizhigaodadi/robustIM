package domain

import (
	"im/ipconfig/pkg/data"
	"log"
	"sort"
	"sync"
)

var (
	Ms *StatManager
)

type StatManager struct {
	endPointMarker *sync.Map
	connCount      int
	eventChan      chan *data.Event
}

func NewStatManager(c chan *data.Event) *StatManager {
	return &StatManager{
		endPointMarker: &sync.Map{},
		connCount:      0,
		eventChan:      c,
	}
}

func (sm *StatManager) IsContains(hostPort string) bool {
	_, ok := sm.endPointMarker.Load(hostPort)
	return ok
}

func (sm *StatManager) Add(event *data.Event) {
	s := &data.Stat{
		ConnectNum:   event.GetConnNum(),
		MessageBytes: event.GetMessageBytes(),
	}
	sm.endPointMarker.Store(event.GetHostPort(), NewSateWindow(DefaultWindowSize, event.GetHost(), int(event.GetPort()), s))
	sm.connCount++
}

func (sm *StatManager) Del(event *data.Event) {
	sm.endPointMarker.Delete(event.GetHostPort())
	sm.connCount--
}

func (sm *StatManager) Update(event *data.Event) {
	stateWindow, ok := sm.endPointMarker.Load(event.GetHostPort())
	if !ok {
		log.Printf("Invalid IpPort:%s\n", event.GetHostPort())
		return
	}
	s := &data.Stat{
		MessageBytes: event.GetMessageBytes(),
		ConnectNum:   event.GetConnNum(),
	}
	stateWindow.(*StateWindow).PushStat(s)

}

func (sm *StatManager) GetIpList() []*EndPoint {
	/*return the best ip list*/
	var endPoints []*EndPoint
	sm.endPointMarker.Range(func(key, value any) bool {
		/*Calculate ranking*/
		sw := value.(*StateWindow)
		endPoints = append(endPoints, &EndPoint{
			Ip:           sw.host,
			Port:         uint16(sw.port),
			staticScore:  sw.GetStaticScore(),
			dynamicScore: sw.GetDynamicScore(),
		})
		return true
	})

	/*Sort*/
	sort.Slice(endPoints, func(i, j int) bool {
		if endPoints[i].dynamicScore != endPoints[j].dynamicScore {
			return endPoints[i].dynamicScore > endPoints[j].dynamicScore
		}
		return endPoints[i].staticScore > endPoints[j].dynamicScore
	})

	return endPoints
}

/*Assert this method will be blocked*/
func (sm *StatManager) RunListening() {
	/*It is processed separately according to the event type*/
	for event := range sm.eventChan {
		t := event.GetEventType()
		if t == data.AddEvent {
			/*Determine whether a new node needs to be added*/
			if !sm.IsContains(event.GetHostPort()) {
				sm.Add(event)
			} else {
				sm.Update(event)
			}
		} else if t == data.DelEvent {
			/*The node goes offline event occurred. Procedure*/
			if !sm.IsContains(event.GetHostPort()) {
				log.Printf("Invalid Ip And Port:%s\n", event.GetHostPort())
			} else {
				sm.Del(event)
			}
		} else { /*UnKnown Err*/
			log.Printf("Invalid EventType For Ip:%s\n", event.GetHostPort())

		}
	}

}
