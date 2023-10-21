package discovery

import "encoding/json"

type EndPointInfo struct {
	Key  string            `json:"Key"`
	Val  string            `json:"Val"`
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
