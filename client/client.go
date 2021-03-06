package main

import (
	"bytes"
	"fmt"
	"github.com/reechou/duobb_access/duobb"
	"github.com/reechou/holmes"
	"github.com/reechou/tao"
	"net"
	"os"
	"time"
)

func main() {
	tao.Register(duobb.DuobbMsgCMD, duobb.DeserializeMessage, nil)

	go taoke()
	//go qunguan()
	
	stop := make(chan struct{})
	<-stop
}

func taoke() {
	//c, err := net.Dial("tcp", "121.40.85.37:7899")
	c, err := net.Dial("tcp", "127.0.0.1:7899")
	if err != nil {
		holmes.Fatal("%v", err)
	}
	
	tcpConnection := tao.NewClientConnection(0, false, c, nil)
	defer tcpConnection.Close()
	
	tcpConnection.SetOnConnectCallback(func(client tao.Connection) bool {
		fmt.Println("On connect")
		return true
	})
	
	tcpConnection.SetOnErrorCallback(func() {
		fmt.Println("%s", "On error")
	})
	
	tcpConnection.SetOnCloseCallback(func(client tao.Connection) {
		fmt.Println("On close")
		os.Exit(0)
	})
	
	tcpConnection.SetOnMessageCallback(func(msg tao.Message, client tao.Connection) {
		fmt.Println(string(msg.(*duobb.DuobbMsg).UserName), string(msg.(*duobb.DuobbMsg).Method), len(msg.(*duobb.DuobbMsg).Msg))
		if string(msg.(*duobb.DuobbMsg).Method) == "DuobbAccountService.LogoutKickOff" {
			fmt.Println("kick off")
			client.Close()
		}
		DecodeMsg(msg.(*duobb.DuobbMsg).Msg)
	})
	
	tcpConnection.Start()
	fmt.Println("start to talk:")
	msg := &duobb.DuobbMsg{
		UserName: []byte("god"),
		Method:   []byte("DuobbAccountService.Login"),
		Msg:      EncodeMsg(),
	}
	tcpConnection.Write(msg)
	time.Sleep(200 * time.Second)
	msg = &duobb.DuobbMsg{
		UserName: []byte("god"),
		Method:   []byte("DuobbAccountService.Heartbeat"),
		Msg:      EncodeHeartbeat(),
	}
	tcpConnection.Write(msg)
	//time.Sleep(15 * time.Second)
	//msg = &duobb.DuobbMsg{
	//	UserName: []byte("god"),
	//	Method:   []byte("DuobbAccountService.GetAllDuobbData"),
	//	Msg:      EncodeMsg(),
	//}
	//tcpConnection.Write(msg)
	//time.Sleep(15 * time.Second)
	//msg = &duobb.DuobbMsg{
	//	UserName: []byte("god"),
	//	Method:   []byte("SelectProductService.GetSpPlanInfoFromUser"),
	//	Msg:      EncodePlanInfoMsg(),
	//}
	//tcpConnection.Write(msg)
	time.Sleep(15 * time.Second)
	msg = &duobb.DuobbMsg{
		UserName: []byte("god"),
		Method:   []byte("DuobbAccountService.Logout"),
		Msg:      EncodeMsg(),
	}
	tcpConnection.Write(msg)
	time.Sleep(5 * time.Second)
	
	tcpConnection.Close()
}

func qunguan() {
	c, err := net.Dial("tcp", "127.0.0.1:7899")
	if err != nil {
		holmes.Fatal("%v", err)
	}
	
	tcpConnection := tao.NewClientConnection(0, false, c, nil)
	defer tcpConnection.Close()
	
	tcpConnection.SetOnConnectCallback(func(client tao.Connection) bool {
		fmt.Println("On connect")
		return true
	})
	
	tcpConnection.SetOnErrorCallback(func() {
		fmt.Println("%s", "On error")
	})
	
	tcpConnection.SetOnCloseCallback(func(client tao.Connection) {
		fmt.Println("On close")
		os.Exit(0)
	})
	
	tcpConnection.SetOnMessageCallback(func(msg tao.Message, client tao.Connection) {
		fmt.Println(string(msg.(*duobb.DuobbMsg).UserName), string(msg.(*duobb.DuobbMsg).Method), len(msg.(*duobb.DuobbMsg).Msg))
		if string(msg.(*duobb.DuobbMsg).Method) == "DuobbAccountService.LogoutKickOff" {
			fmt.Println("kick off")
			client.Close()
		}
		DecodeMsg(msg.(*duobb.DuobbMsg).Msg)
	})
	
	tcpConnection.Start()
	fmt.Println("start to talk:")
	msg := &duobb.DuobbMsg{
		AppId:    1,
		UserName: []byte("god"),
		Method:   []byte("DuobbAccountService.Login"),
		Msg:      EncodeMsg(),
	}
	tcpConnection.Write(msg)
	time.Sleep(5 * time.Second)
	msg = &duobb.DuobbMsg{
		AppId:    1,
		UserName: []byte("god"),
		Method:   []byte("DuobbAccountService.Heartbeat"),
		Msg:      EncodeHeartbeat(),
	}
	tcpConnection.Write(msg)
	time.Sleep(15 * time.Second)
	msg = &duobb.DuobbMsg{
		AppId:    1,
		UserName: []byte("god"),
		Method:   []byte("DuobbAccountService.Logout"),
		Msg:      EncodeMsg(),
	}
	tcpConnection.Write(msg)
	time.Sleep(5 * time.Second)
	
	tcpConnection.Close()
}

func EncodeMsg() []byte {
	s := &duobb.Security{}
	secretKey := s.Md5Of32(s.Md5Of32([]byte("123456god")))
	msgEncode1 := s.Base64Encode([]byte(`{"user": "god"}`))
	msgEncode1 = append(secretKey, msgEncode1...)
	result, _ := s.GzipEncode(msgEncode1)
	return result
}

func EncodeHeartbeat() []byte {
	s := &duobb.Security{}
	secretKey := s.Md5Of32(s.Md5Of32([]byte("123456god")))
	msgEncode1 := s.Base64Encode([]byte(`{"user":"god","lastPushMsgTime":0}`))
	msgEncode1 = append(secretKey, msgEncode1...)
	result, _ := s.GzipEncode(msgEncode1)
	return result
}

func EncodePlanMsg() []byte {
	s := &duobb.Security{}
	secretKey := s.Md5Of32(s.Md5Of32([]byte("123456god")))
	msgEncode1 := s.Base64Encode([]byte(`{"user":"god","offset":0,"num":100}`))
	msgEncode1 = append(secretKey, msgEncode1...)
	result, _ := s.GzipEncode(msgEncode1)
	return result
}

func EncodePlanInfoMsg() []byte {
	s := &duobb.Security{}
	secretKey := s.Md5Of32(s.Md5Of32([]byte("123456god")))
	msgEncode1 := s.Base64Encode([]byte(`{"user":"god","planId":10}`))
	msgEncode1 = append(secretKey, msgEncode1...)
	result, _ := s.GzipEncode(msgEncode1)
	return result
}

func DecodeMsg(in []byte) {
	s := &duobb.Security{}
	secretKey := s.Md5Of32(s.Md5Of32([]byte("123456god")))
	msgDecode1, err := s.GzipDecode(in)
	if err != nil {
		return
	}
	msgDecode2 := bytes.Replace(msgDecode1, secretKey, []byte(""), -1)

	result, _ := s.Base64Decode(string(msgDecode2))
	fmt.Println("result:", string(result))
}
