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
	host string
	serv_channel_out chan string
	serv_channel_in chan common.Server_storage_data
}



func New(host string, serv_channel_out chan string, serv_channel_in chan common.Server_storage_data) server{
	return server{
		host: host,
		serv_channel_out: serv_channel_out,
		serv_channel_in: serv_channel_in,
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
    serv.serv_channel_out <- params.Get("q")
	out_buf := <- serv.serv_channel_in
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
	http.ListenAndServe(":" + serv.host, nil)
}