package subscriber

import (
	"fmt"
	"log"
	"server-issuing-orders/common"
	"time"

	"github.com/nats-io/stan.go"
)

type subscriber struct{
	ClientID string
	ClusterID string
	ServerURL string
	ChannelName string
}

func New(ClientID, ClusterID, ChannelName string, ServerPort int) subscriber{
	return subscriber{
		ClientID: ClientID,
		ClusterID: ClusterID,
		ChannelName: ChannelName,
		ServerURL: fmt.Sprintf("nats://localhost:%d", ServerPort),
	}
}

func (sub *subscriber)DataFromServer() (chan []byte, error) {
	sc, err := stan.Connect(sub.ClusterID, sub.ClientID, stan.NatsURL(sub.ServerURL))
	if err != nil{
		return nil, common.Wrap("Unable to connect to NATS streaming server", err)
	}
	log.Println("Nats streaming is listening")
	ch := make(chan []byte)
	_, err = sc.Subscribe("orders",
  	func(m *stan.Msg) {
		log.Println("Subscriver recieve data from Nats streaming")
		ch <- m.Data
	}, stan.DeliverAllAvailable(), stan.AckWait(20*time.Second))
	if err != nil{
		return nil, common.Wrap("Unable to connect to NATS streaming server", err)
	}
	return ch, nil
}