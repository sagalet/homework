package common

import (
	"encoding/json"
	"io/ioutil"
)

type Server struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func GetServerInfoFromFile(file string) (*Server, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	s := &Server{}
	err = json.Unmarshal(data, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
