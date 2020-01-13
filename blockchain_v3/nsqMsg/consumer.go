package nsqMsg

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
)

var Consumer *nsq.Consumer

const (
	PropTimeout = "propTimeout" // 道具过期通知
)

func init() {
	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	var err error
	Consumer, err = nsq.NewConsumer(topicName, ChannelName, config)
	if err != nil {
		panic(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	Consumer.AddHandler(&myMessageHandler{})

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = Consumer.ConnectToNSQLookupd(nsqUrl + ":4161")
	if err != nil {
		panic(err)
	}

}

type myMessageHandler struct {}

// HandleMessage implements the Handler interface.
func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		return nil
	}

	err := processMessage(m.Body)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}



func processMessage(msg []byte) error {
	log.Printf("RawMsg: %s", msg)
	handle := Handle{}
	err := json.Unmarshal(msg, &handle)
	if err != nil {
		return err
	}

	switch handle.Type {
	case PropTimeout:
		err = handlePropTimeout(handle.Js)
	default:
		log.Printf("Unknown Handle Type %s", handle.Type)
		return nil
	}

	if err != nil {
		return err
	}

	log.Println("Handle success.")
	return nil
}