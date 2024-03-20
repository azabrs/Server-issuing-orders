package subscriber

import (
	"fmt"
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

func (sub *subscriber)DataFromServer() error {
	sc, err := stan.Connect(sub.Cluster_ID, sub.Client_ID, stan.NatsURL(sub.Server_URL))
	if err != nil{
		return err
	}
	_, err = sc.Subscribe("orders",
  	func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	}, stan.DeliverAllAvailable(), stan.AckWait(20*time.Second))
	if err != nil{
		return err
	}
	return nil
}