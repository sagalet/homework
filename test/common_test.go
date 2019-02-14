package test

import (
	"flag"
	"io/ioutil"
	"main/common"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInsertData(t *testing.T) {
	common.InitRedis("127.0.0.1:6379")

	// try insert 60 requests
	for i := 1; i <= 60; i++ {
		n, err := common.InsertRequest("192.168.0.0")
		if err != nil {
			t.Errorf("get error %q when insert %d", err.Error(), i)
		} else if n != i {
			t.Errorf("request number is not match %d but get %d", i, n)
		}
	}

	// 61 and error should be returned
	n, err := common.InsertRequest("192.168.0.0")
	if n != 61 || err == nil {
		t.Errorf("the error should be return when insert 61'th requests but get %d err %q", n, err.Error())
	}

	common.DeleteIP("192.168.0.0")
}

func insertHTTPRequestInternel(t *testing.T, addr string) {

	// try insert 60 requests
	for i := 1; i <= 60; i++ {
		resp, err := http.Get(addr)
		if err != nil {
			t.Errorf("Connect server failed %q", err.Error())
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		n, err := strconv.Atoi(string(body[:]))
		if err != nil {
			t.Errorf("get error %q when insert %d", err.Error(), i)
		} else if n != i {
			t.Errorf("request number is not match %d but get %d", i, n)
		}
		resp.Body.Close()
	}

	// error should be returned
	resp, err := http.Get(addr)
	if err != nil {
		t.Errorf("Connect server failed %q", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	s := string(body[:])
	if strings.Compare(s, "error") != 0 {
		t.Errorf("Error should be returned if 61'th request sent but get %s", s)
	}

	resp.Body.Close()
}

func TestInsertHTTPRequest(t *testing.T) {
	port := flag.String("port", "10034", "server port")
	flag.Parse()
	addr := "http://:" + *port
	insertHTTPRequestInternel(t, addr)
	// sleep 61 second and test again ti check timeout
	time.Sleep(time.Minute + time.Second)
	insertHTTPRequestInternel(t, addr)
}
