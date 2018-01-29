package main

import "testing"

func Test_Binary(t *testing.T) {
	pool := GetBufferPool()
	index := 0
	for i := 0; i < 10; i++ {
		index = pool.Add(NewBufferWriter(i))
		t.Log(i, "=", index)
	}
	t.Log(pool.Len())

	//for i := 0; i < 10; i++ {
	//	buffer := pool.Get(i)
	//	if buffer != nil {
	//		t.Log(i, " cap ", buffer.Cap())
	//	}
	//}

	for i := 0; i <= 11; i++ {
		index = pool.find(i)
		t.Log(i, " index ", index)
	}

}
