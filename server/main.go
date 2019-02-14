package main

import (
	"flag"
	"log"
	"main/common"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Handler struct{}

func (this Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.RemoteAddr, ":")
	log.Println("get request from ", addr[0])
	index, err := common.InsertRequest(addr[0])
	if err != nil {
		if index > 60 {
			w.Write([]byte("error"))
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
		w.Write([]byte(strconv.Itoa(index)))
	}
}

func main() {

	port := flag.String("port", "10034", "server port")
	conf := flag.String("conf", "redis.conf", "redis server info")

	flag.Parse()

	h := Handler{}
	s := &http.Server{
		Addr:         ":" + *port,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	r, err := common.GetServerInfoFromFile(*conf)

	if err != nil {
		log.Println("err: ", err)
		return
	}

	common.InitRedis(r.IP + ":" + r.Port)

	log.Println("start server")

	go func(s *http.Server) {
		log.Println("Press 'Enter' to stop server")
		b := make([]byte, 1)
		os.Stdin.Read(b)
		s.Close()
	}(s)

	log.Fatal(s.ListenAndServe())

}
