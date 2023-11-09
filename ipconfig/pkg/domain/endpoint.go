package domain

import "fmt"

/*Return the model to the client*/
type EndPoint struct {
	Ip           string `json:"ip"`
	Port         uint16 `json:"port"`
	staticScore  float64
	dynamicScore float64
}

func (e *EndPoint) GetHostStr() string {
	return fmt.Sprintf("%s:%d", e.Ip, e.Port)
}
