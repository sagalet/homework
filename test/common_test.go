package test

import (
	"flag"
	"io/ioutil"
	"log"
	"main/common"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	mRedisIp    = flag.String("redisip", "127.0.0.1", "redis ip")
	mRedisPort  = flag.String("redisport", "6379", "redis port")
	mServerPort = flag.String("port", "10034", "server port")
)

func TestInsertData(t *testing.T) {
	flag.Parse()

	addr := *mRedisIp + ":" + *mRedisPort
	t.Logf("Connect server %s", addr)
	common.InitRedis(addr)

	defer common.DeleteIP("192.168.0.0")
	// try insert 60 requests
	for i := 1; i <= 60; i++ {
		n, err := common.InsertRequest("192.168.0.0")
		if err != nil {
			t.Errorf("get error %q when insert %d", err.Error(), i)
			return
		} else if n != i {
			t.Errorf("request number is not match %d but get %d", i, n)
		}
	}

	// 61 and error should be returned
	n, err := common.InsertRequest("192.168.0.0")
	if n != 61 || err == nil {
		t.Errorf("the error should be return when insert 61'th requests but get %d err %q", n, err.Error())
	}
}

func insertHTTPRequestInternel(t *testing.T, addr string) error {

	// try insert 60 requests
	for i := 1; i <= 60; i++ {
		resp, err := http.Get(addr)
		if err != nil {
			t.Errorf("Connect server failed %q", err.Error())
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		n, err := strconv.Atoi(string(body[:]))
		if err != nil {
			t.Errorf("get error %q when insert %d", err.Error(), i)
			return err
		} else if n != i {
			t.Errorf("request number is not match %d but get %d", i, n)
		}
		resp.Body.Close()
	}

	// "error" should be returned
	resp, err := http.Get(addr)
	if err != nil {
		t.Errorf("Connect server failed %q", err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	s := string(body[:])
	if strings.Compare(s, "error") != 0 {
		t.Errorf("Error should be returned if 61'th request sent but get %s", s)
	}

	resp.Body.Close()
	return nil
}

func TestInsertHTTPRequest(t *testing.T) {
	flag.Parse()
	addr := "http://:" + *mServerPort
	err := insertHTTPRequestInternel(t, addr)
	if err != nil {
		return
	}
	log.Println("Sleep 61 second for expired and test again")
	time.Sleep(time.Minute + time.Second)
	err = insertHTTPRequestInternel(t, addr)
}
