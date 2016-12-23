package duobb

import (
	"bytes"
	"encoding/binary"

	"github.com/reechou/holmes"
	"github.com/reechou/tao"
)

const (
	DuobbMsgCMD int32 = 0x01
)

type DuobbMsg struct {
	UserName []byte
	Method   []byte
	Msg      []byte
}

func (self *DuobbMsg) MessageNumber() int32 {
	return DuobbMsgCMD
}

func (self *DuobbMsg) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(len(self.UserName)))
	buf.Write(self.UserName)
	binary.Write(buf, binary.LittleEndian, int32(len(self.Method)))
	buf.Write(self.Method)
	binary.Write(buf, binary.LittleEndian, int32(len(self.Msg)))
	buf.Write(self.Msg)

	return buf.Bytes(), nil
}

// 4bytes-len | username | 4bytes-len | method | 4bytes-len | msg
func DeserializeMessage(data []byte) (tao.Message, error) {
	if data == nil {
		return nil, ErrorNilData
	}
	dataLen := uint32(len(data))

	buffer := bytes.NewBuffer(data)

	var len uint32
	binary.Read(buffer, binary.LittleEndian, &len)
	dataLen -= 4
	if len > dataLen {
		return nil, ErrorIllegalData
	}
	userNameBytes := make([]byte, len)
	binary.Read(buffer, binary.LittleEndian, userNameBytes)
	dataLen -= len

	binary.Read(buffer, binary.LittleEndian, &len)
	dataLen -= 4
	if len > dataLen {
		return nil, ErrorIllegalData
	}
	methodBytes := make([]byte, len)
	binary.Read(buffer, binary.LittleEndian, methodBytes)
	dataLen -= len

	binary.Read(buffer, binary.LittleEndian, &len)
	dataLen -= 4
	if len > dataLen {
		return nil, ErrorIllegalData
	}
	msgBytes := make([]byte, len)
	binary.Read(buffer, binary.LittleEndian, msgBytes)

	msg := &DuobbMsg{
		UserName: userNameBytes,
		Method:   methodBytes,
		Msg:      msgBytes,
	}
	return msg, nil
}

func ProcessDuobbMessage(ctx tao.Context, conn tao.Connection) {
	msg := ctx.Message().(*DuobbMsg)
	if msg == nil {
		holmes.Error("error duobb msg: %v.", msg)
		return
	}
	holmes.Debug(string(msg.UserName))
	holmes.Debug(string(msg.Method))
	holmes.Debug(string(msg.Msg))
}
