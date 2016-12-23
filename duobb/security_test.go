package duobb

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	s := &Security{}
	r := s.Md5Of32([]byte("xxxxxxxx"))
	fmt.Printf("%s \n", string(r))
	r = s.Md5Of32(r)
	fmt.Printf("%s \n", string(r))
}

func TestBase64(t *testing.T) {
	s := &Security{}
	r := s.Base64Encode([]byte("xxxxxxxx"))
	fmt.Printf("%s \n", r)
	r2, _ := s.Base64Decode(r)
	fmt.Printf("%s \n", string(r2))
}
