package discovery

import "encoding/json"

type EndPointInfo struct {
	Ip   string            `json:"Key"`
	Port string            `json:"Port"`
	Meta map[string]string `json:"Meta"`
}

func (e *EndPointInfo) Marshal() []byte {
	serial, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return serial
}

func UnMarshal(serial []byte) *EndPointInfo {
	var e = &EndPointInfo{}
	err := json.Unmarshal(serial, e)
	if err != nil {
		return nil
	}
	return e
}
