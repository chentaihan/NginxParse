package util

import (
	"testing"
)

func Test_Write(t *testing.T) {
	t.Log("Test_Write------------------------------------")
	writer := NewBufferWriter(0)
	writer.Write([]byte("0123456789"))
	writer.Write([]byte("qazwsxedcrfv"))
	t.Log(writer.ToString())

	buf := writer.buffer[0:0]
	writer.buffer = buf

	buf = append(buf, "XXXXXXXX"...)
	t.Log(string(buf))
	t.Log(writer.ToString())
	t.Log(cap(writer.buffer))
}

func Test_ToString(t *testing.T){
	t.Log("Test_ToString------------------------------------")
	writer := NewBufferWriter(0)
	t.Log(writer.ToString())
	writer.Write([]byte("0123456789"))
	t.Log(writer.ToString())
	writer.Clear()
	writer.Write([]byte("zxc"))
	t.Log(writer.ToString())
}

func Test_MoveNext(t *testing.T){
	t.Log("Test_MoveNext------------------------------------")
	writer := NewBufferWriter(0)

	writer.Write([]byte("0123456789\n"))
	writer.Write([]byte("0123456789\n"))
	writer.Write([]byte("0123456789\n"))
	writer.Write([]byte("0123456789\n"))

	writer.Reset()
	for writer.MoveNext(){
		t.Log(writer.Current())
	}
	writer.Write([]byte("01234567"))

	writer.Reset()
	for writer.MoveNext(){
		t.Log(writer.Current())
	}
}