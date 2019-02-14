package main

import (
	"flag"
	"io/ioutil"
	"log"
	"main/common"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	delay := flag.Int("delay", 2000, "delay between each request")
	conf := flag.String("conf", "server.conf", "server info")
	count := flag.Int("count", 120, "the number of requests")
	flag.Parse()

	// now try to send 120 requests with delay
	s, err := common.GetServerInfoFromFile(*conf)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	addr := "http://" + s.IP + ":" + s.Port
	for i := 1; i <= *count; i++ {
		resp, err := http.Get(addr)
		if err != nil {
			log.Println("Connect server failed %q", err.Error())
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("get error %q when insert %d", err.Error(), i)
			continue
		}
		log.Printf("get(%d) : %s\n", i, string(body[:]))

		resp.Body.Close()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(*delay)))
	}
}
