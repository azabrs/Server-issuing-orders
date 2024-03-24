package main

import (
	"log"
	server "server-issuing-orders/Server"
	storage "server-issuing-orders/Storage"
	subscriber "server-issuing-orders/Subscriber"
)


func main(){
	chan_server_storage := make(chan string)
	sub := subscriber.New("listener", "test-cluster", "orders", 4222)
	ch, err := sub.DataFromServer()
	if err != nil{
		log.Fatal(err)
	}
	stor, err := storage.New(ch, chan_server_storage, "apuha", "12345678", "wb", "wb_table")
	if err != nil{
		log.Fatal(err)
	}
	if err = stor.Handler(); err != nil{
		log.Fatal(err)
	}
	serv := server.New("8080", chan_server_storage, stor.Serv_channel_out)
	serv.StartServer()
}