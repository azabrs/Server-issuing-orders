package subscriber

import (
	"fmt"
	"log"
	"server-issuing-orders/common"
	"time"

	"github.com/nats-io/stan.go"
)

type subscriber struct{
	Client_ID string
	Cluster_ID string
	Server_URL string
	Channel_name string
}

func New(Client_ID, Cluster_ID, Channel_name string, Server_port int) subscriber{
	return subscriber{
		Client_ID: Client_ID,
		Cluster_ID: Cluster_ID,
		Channel_name: Channel_name,
		Server_URL: fmt.Sprintf("nats://localhost:%d", Server_port),
	}
}

func (sub *subscriber)DataFromServer() (chan []byte, error) {
	sc, err := stan.Connect(sub.Cluster_ID, sub.Client_ID, stan.NatsURL(sub.Server_URL))
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