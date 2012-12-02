package bencode

import (
	"bytes"
	"testing"
)

type encoderTest struct {
	v   interface{}
	ret string
}

var encTests = []encoderTest{
	{0, "i0e"},
	{10, "i10e"},
	{-10, "i-10e"},
	{"test", "4:test"},
	{[]byte{'a', 'b', 'c'}, "3:abc"},
}

func TestEncoder(t *testing.T) {
	buf := new(bytes.Buffer)
	for i, test := range encTests {
		buf.Reset()
		e := NewEncoder(buf)

		err := e.Encode(test.v)
		if err != nil {
			t.Errorf("Test %d: returned error `%s`.", i, err.Error())
		}
		ret := buf.String()

		if ret != test.ret {
			t.Errorf("Test %d: expected `%s`, got `%s`.", i, test.ret, ret)
		}
	}
}
