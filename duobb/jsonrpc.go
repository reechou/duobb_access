package duobb

import (
	"bytes"
	"net/http"
	//"io/ioutil"

	"github.com/gorilla/rpc/json"
	"github.com/reechou/duobb_proto"
	"github.com/reechou/holmes"
)

type JsonRpc struct {
	client *http.Client
}

func NewJsonRpc() *JsonRpc {
	return &JsonRpc{
		client: &http.Client{},
	}
}

func (self *JsonRpc) Call(host, method string, request interface{}) (interface{}, error) {
	url := "http://" + host + "/rpc"
	message, err := json.EncodeClientRequest(method, request)
	if err != nil {
		holmes.Error("json encode client request error: %v", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
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
	//rspBody, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	holmes.Error("ioutil ReadAll error: %v", err)
	//	return nil, err
	//}
	var result duobb_proto.Response
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		holmes.Error("json decode client response error: %v", err)
		return nil, err
	}
	return result, nil
}
