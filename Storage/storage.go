package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"server-issuing-orders/common"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)


type storage struct{
	Cache map[string] common.Order
	PostgresDataSource string
	SubChannel chan []byte
	ServChannelIn chan string
	ServChannelOut chan common.ServerStorageData
	TableName string
	db *sql.DB

}

func New(SubChannel chan []byte, ServChannelIn chan string,  user, password, dbname, TableName string) (storage, error){
	stor := storage{
		SubChannel: SubChannel,
		ServChannelIn: ServChannelIn,
		PostgresDataSource: fmt.Sprintf("dbname=%s user=%s password=%s sslmode = disable", dbname, user, password),
		TableName: TableName,
		Cache: make(map[string]common.Order),
	}
	ServChannelOut := make(chan common.ServerStorageData)
	stor.ServChannelOut = ServChannelOut
	if err := stor.InitCache(); err != nil{
		return storage{}, common.Wrap("Unable to initialize cache from postgres", err)
	}
	return stor, nil
}

func(s *storage) InitCache() error{
	db, err := sql.Open("postgres", s.PostgresDataSource)
	if err != nil{
		common.Wrap("Сan't connect to the database", err)
	}
	quer := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (OrderUID VARCHAR(30), data VARCHAR(20000));", s.TableName)
	if _, err := db.Exec(quer); err != nil{
		common.Wrap("Problem with the table", err)
	}
	s.db = db
	quer = fmt.Sprintf("SELECT * FROM %s", s.TableName)
	rows, err := db.Query(quer)
	if err != nil{
		return common.Wrap("Data could not be retrieved from the table", err)
	}
	for rows.Next(){
		order := new(common.Order)
		var id string
		var data string
		var dataByte []byte
		if err := rows.Scan(&id, &data); err != nil{
			return common.Wrap("The data received from the table could not be converted", err)
		}
		data = strings.Trim(data, "[]")
		temp := strings.Split(data, " ")
		for _, val := range temp{
			buf_temp, _ := strconv.Atoi(val)
			dataByte = append(dataByte, byte(buf_temp))
		}
		if err := json.Unmarshal(dataByte, order); err != nil{
			return common.Wrap("The data received from the table could not be converted", err)
		}
		s.Cache[id] = *order
	}
	return nil
}

func (s *storage) Handler() error{
	go func(){
		for{
			select{
			case buf := <- s.SubChannel:
				log.Println("Storage receive data from Nats streaming")
				var order common.Order
				if err := json.Unmarshal([]byte(buf), &order); err != nil{
					log.Fatal(err)
				}
				if _, ok := s.Cache[order.OrderUID]; ok {
					log.Printf("Order with %s UID is already in the database, order skipped", order.OrderUID)
				} else{
					s.Cache[order.OrderUID] = order
					temp, _ := json.Marshal(order)
					q := fmt.Sprintf("INSERT INTO %s VALUES('%s', '%v');", s.TableName, order.OrderUID, temp)
					if _, err := s.db.Exec(q); err != nil{
						log.Fatal(err)
					}
					log.Println("Data has been successfully added to the postgres")
				}
				
			case buf := <- s.ServChannelIn:
				log.Println("Storage received a request from the server")
				if val, ok := s.Cache[buf]; ok {
					buf_out, _ := json.Marshal(val)
					s.ServChannelOut <- common.ServerStorageData{Exist: true, Data: buf_out}
				} else{
					s.ServChannelOut <- common.ServerStorageData{Exist: false, Data: nil}
				}

			}
		}
	} ()

	return nil
}