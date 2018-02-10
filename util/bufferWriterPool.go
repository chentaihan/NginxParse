package util

/**
排序slice存放BufferWriter
Add：添加元素
所有被使用的元素都会从slice中移除，如果需要回收请调Add
*/

var bufferWriterPool *BufferWriterPool

type BufferWriterPool struct {
	pool []*BufferWriter
}

func NewBufferWriter(cap int) *BufferWriter {
	return GetBufferPool().NewBuffer(cap)
}

func GetBufferPool() *BufferWriterPool {
	if bufferWriterPool == nil {
		bufferWriterPool = &BufferWriterPool{}
		for i := 0; i < 4; i++{
			bufferWriterPool.Add(newBufferWriter(0))
		}
	}
	return bufferWriterPool
}

func (pool *BufferWriterPool) Add(writer *BufferWriter) int {
	writer.Clear()
	index := pool.find(writer.Cap())
	if index < pool.Len()-1 {
		after := pool.pool[index:]
		pool.pool = append(pool.pool[0:index], writer)
		pool.pool = append(pool.pool, after...)
	} else {
		pool.pool = append(pool.pool, writer)
	}
	return index
}

//按照容量二分查找
func (pool *BufferWriterPool) find(cap int) int {
	start := 0
	end := len(pool.pool) - 1
	index := 0
	isBig := 0
	for start <= end {
		index = start + (end-start)/2
		itemCap := pool.pool[index].Cap()
		if itemCap == cap {
			isBig = 0
			break
		} else if itemCap < cap {
			start = index + 1
			isBig = 1
		} else {
			end = index - 1
			isBig = 0
		}
	}

	return index + isBig
}

func (pool *BufferWriterPool) Get(index int) *BufferWriter {
	poolLen := pool.Len()
	if poolLen == 0 {
		return newBufferWriter(0)
	}
	if index < 0 || index >= poolLen {
		return pool.Remove(poolLen / 2)
	}
	return pool.Remove(index)
}

func (pool *BufferWriterPool) Remove(index int) *BufferWriter {
	if index < 0 || index >= pool.Len() {
		return nil
	}

	ret := pool.pool[index]
	if index == pool.Len()-1 {
		pool.pool = pool.pool[0:index]
	} else {
		pool.pool = append(pool.pool[0:index], pool.pool[index+1:]...)
	}

	return ret
}

func (pool *BufferWriterPool) NewBuffer(cap int) *BufferWriter {
	index := pool.find(cap)
	if index >= pool.Len() {
		return newBufferWriter(cap)
	}
	return pool.Remove(index)
}

func (pool *BufferWriterPool) Len() int {
	return len(pool.pool)
}

func (pool *BufferWriterPool) Cap() int {
	return cap(pool.pool)
}
