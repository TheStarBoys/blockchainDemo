package nsqMsg

import (
	"github.com/nsqio/go-nsq"
	"log"
)

var Producer *nsq.Producer
const (
	nsqUrl         = "172.27.16.17"
	topicName      = "topic"
	ChannelName    = "channel"
)

func init() {
	config := nsq.NewConfig()
	var err error
	Producer, err = nsq.NewProducer(nsqUrl + ":4150", config)
	if err != nil {
		log.Panic(err)
	}
}

func NsqPub(msg []byte) error {
	return Producer.Publish(topicName, msg)
}