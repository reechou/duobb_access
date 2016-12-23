package duobb

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"

	"github.com/reechou/holmes"
)

type Security struct{}

func (self *Security) Md5Of32(src []byte) []byte {
	hash := md5.New()
	hash.Write(src)
	cipherText2 := hash.Sum(nil)
	hexText := make([]byte, 32)
	hex.Encode(hexText, cipherText2)
	return hexText
}

func (self *Security) Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)
	return buf
}

func (self *Security) Base64Decode(src string) ([]byte, error) {
	r, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		holmes.Error("base64 decode src[%s] error: %v", src, err)
		return nil, ErrorBase64Decode
	}
	return r, nil
}

func (self *Security) GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}
	return buffer.Bytes(), nil
}

func (self *Security) GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
