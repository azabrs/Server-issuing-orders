package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"server-issuing-orders/common"
	_ "github.com/lib/pq"
)


type storage struct{
	Cache map[string] common.Order
	PostgresDataSource string
	sub_channel chan []byte
	serv_channel_in chan string
	Serv_channel_out chan common.Server_storage_data
	table_name string
	db *sql.DB

}

func New(sub_channel chan []byte, serv_channel_in chan string,  user, password, dbname, table_name string) (storage, error){
	stor := storage{
		sub_channel: sub_channel,
		serv_channel_in: serv_channel_in,
		PostgresDataSource: fmt.Sprintf("dbname=%s user=%s password=%s sslmode = disable", dbname, user, password),
		table_name: table_name,
	}
	Serv_channel_out := make(chan common.Server_storage_data)
	stor.Serv_channel_out = Serv_channel_out
	if err := stor.InitCache(); err != nil{
		return storage{}, err
	}
	return stor, nil
}

func(s *storage) InitCache() error{
	db, err := sql.Open("postgres", s.PostgresDataSource)
	if err != nil{
		return err
	}
	quer := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (OrderUID VARCHAR(30), data VARCHAR(600));", s.table_name)
	if _, err := db.Exec(quer); err != nil{
		return err
	}
	s.db = db
	quer = fmt.Sprintf("SELECT * FROM %s", s.table_name)
	rows, err := db.Query(quer)
	if err != nil{
		return err
	}
	for rows.Next(){
		var order common.Order
		rows.Scan(order)
		//err = json.Unmarshal([]byte(val), &order)
		//if err != nil{
		//	return err
		//}
		s.Cache[order.OrderUID] = order
	}
	return nil
}

func (s *storage) Handler() error{
	go func(){
		for{
			select{
			case buf := <- s.sub_channel:
				var order common.Order
				if err := json.Unmarshal([]byte(buf), &order); err != nil{
					log.Fatal(err)
				}
				if _, ok := s.Cache[order.OrderUID]; ok {
					log.Printf("Order with %s UID is already in the database, order skipped", order.OrderUID)
				} else{
					s.Cache[order.OrderUID] = order
					if _, err := s.db.Exec("INSERT INTO $1 VALUES($2, $3);", s.table_name, order.OrderUID, order); err != nil{
						log.Fatal(err)
					}
				}
				
			case buf := <- s.serv_channel_in:
				if val, ok := s.Cache[buf]; ok {
					buf_out, _ := json.Marshal(val)
					s.Serv_channel_out <- common.Server_storage_data{Exist: true, Data: buf_out}
				} else{
					s.Serv_channel_out <- common.Server_storage_data{Exist: false, Data: nil}
				}

			}
		}
	} ()

	return nil
}