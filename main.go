package main

import (
	"log"
	server "server-issuing-orders/Server"
	storage "server-issuing-orders/Storage"
	subscriber "server-issuing-orders/Subscriber"
	"flag"
)

func mustToken() []string{
	var res []string
	s1 := flag.String("name", "", "account name")
	s2 := flag.String("password", "", "password for account")
	flag.Parse()
	res = append(res, *s1)
	res = append(res, *s2)
	if res[0] == "" || res[1] == ""{
		log.Fatal("Account name and password must set")
	}
	return res
}

func main(){
	flags := mustToken()
	chan_server_storage := make(chan string)
	sub := subscriber.New("listener", "test-cluster", "orders", 4222)
	ch, err := sub.DataFromServer()
	if err != nil{
		log.Fatal(err)
	}
	stor, err := storage.New(ch, chan_server_storage, flags[0], flags[1], "wb", "wb_table")
	if err != nil{
		log.Fatal(err)
	}
	if err = stor.Handler(); err != nil{
		log.Fatal(err)
	}
	serv := server.New("8080", chan_server_storage, stor.Serv_channel_out)
	serv.StartServer()
}