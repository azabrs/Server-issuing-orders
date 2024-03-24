package server

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"server-issuing-orders/common"
)


var tpl = template.Must(template.ParseFiles("server.html"))


type server struct{
	Host string
	ServChannelOut chan string
	ServChannelIn chan common.ServerStorageData
}



func New(host string, ServChannelOut chan string, ServChannelIn chan common.ServerStorageData) server{
	return server{
		Host: host,
		ServChannelOut : ServChannelOut,
		ServChannelIn : ServChannelIn,
	}

}



func ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tpl.Execute(w, nil)
}

func (serv *server)OrderUIDHandel(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    params := u.Query()
    serv.ServChannelOut <- params.Get("q")
	out_buf := <- serv.ServChannelIn
	if !out_buf.Exist{
		w.Write([]byte("No such orderUID was found"))
		log.Println("No such orderUID was found")
	} else{
		log.Printf("Data found with orderUID %s", params.Get("q"))
		w.Write(out_buf.Data)
	}

}


func(serv *server) StartServer(){
	http.HandleFunc("/orderUID", serv.OrderUIDHandel)
    http.HandleFunc("/", ServeHTTP)
	log.Println("Server is listening")
	http.ListenAndServe(":" + serv.Host, nil)
}