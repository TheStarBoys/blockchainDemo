package nsqMsg

import (
	"encoding/json"
	"fmt"
)

type MsgPropTimeout struct {
	PlayerId     string `json:"player_id"`
	MarketItemId string `json:"market_item_id"`
}

func (m *MsgPropTimeout) Pub() error {
	js, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("MsgPropTimeout json: %v", err)
	}

	handle := Handle{
		Type: PropTimeout,
		Js:   js,
	}

	return NsqPub(handle.ToJson())
}