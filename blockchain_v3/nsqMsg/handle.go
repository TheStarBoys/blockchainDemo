package nsqMsg

import (
	"encoding/json"
)

//统一接受参数
type Handle struct {
	Type string `json:"type"` //nsq类型，用于判断
	Js   []byte `json:"js"`   //对象json
}

func (h *Handle) ToJson() []byte {
	js, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}

	return js
}

func handlePropTimeout(msg []byte) error {
	m := MsgPropTimeout{}
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return err
	}

	// goworld do something
	return nil
}