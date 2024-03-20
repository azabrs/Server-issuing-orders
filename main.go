package main

import (
	"log"
	subscriber "server-issuing-orders/Subscriber"
	"time"
)


func main(){
	sub := subscriber.New("listener", "test-cluster", "orders", 4222)
	err := sub.DataFromServer()
	if err != nil{
		log.Fatal(err)
	}
	time.Sleep(time.Second * 10)
}