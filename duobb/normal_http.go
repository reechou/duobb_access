package duobb

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/reechou/holmes"
)

type NormalHttp struct {
	client *http.Client
}

func NewNormalHttp() *NormalHttp {
	return &NormalHttp{
		client: &http.Client{},
	}
}

func (self *NormalHttp) Call(host, uri string, request []byte) ([]byte, error) {
	url := "http://" + host + uri
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	return rspBody, nil
}
