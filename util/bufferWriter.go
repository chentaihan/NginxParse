package util

import "unsafe"

const (
	BUFFER_SIZE = 1024
	CUT_CHAR    = '\n'
)

type LineInfo struct {
	Start int
	End   int
}

type BufferWriter struct {
	buffer  []byte
	line    LineInfo
	cutChar byte //分割字符
}

func newBufferWriter(cap int) *BufferWriter {
	writer := &BufferWriter{}
	if cap <= 0 {
		cap = BUFFER_SIZE
	}
	writer.buffer = make([]byte, 0, cap)
	writer.cutChar = CUT_CHAR
	return writer
}

func (writer *BufferWriter) Write(p []byte) (n int, err error) {
	writer.buffer = append(writer.buffer, p...)
	n = len(p)
	err = nil
	return n, err
}

func (writer *BufferWriter) WriteChar(p byte) (n int, err error) {
	writer.buffer = append(writer.buffer, p)
	err = nil
	return 1, err
}

func (writer *BufferWriter) WriteString(str string) (n int, err error) {
	p := []byte(str)
	writer.buffer = append(writer.buffer, p...)
	n = len(p)
	err = nil
	return n, err
}

func (writer *BufferWriter) Clear() {
	writer.buffer = writer.buffer[0:0]
	writer.Reset()
}

func (writer *BufferWriter) GetBuffer() []byte {
	return writer.buffer
}

/**
在多协程的情况，不能这样写，需要复制
 */
func (writer *BufferWriter) ToString() string {
	buf := writer.GetBuffer()
	return *(*string)(unsafe.Pointer(&buf))
}

func (writer *BufferWriter) MoveNext() bool {
	size := len(writer.buffer)
	if writer.line.End >= size {
		return false
	}
	writer.line.Start = writer.line.End
	i := writer.line.Start
	for ; i < size; i++ {
		if writer.buffer[i] == writer.cutChar {
			writer.line.End = i + 1
			break
		}
	}
	if i >= size {
		writer.line.End = i
	}
	return true
}

func (writer *BufferWriter) IsEnd() bool {
	if writer.line.End >= len(writer.buffer) {
		return true
	}
	return false
}

func (writer *BufferWriter) Current() string {
	current := writer.buffer[writer.line.Start:writer.line.End]
	return *(*string)(unsafe.Pointer(&current))
}

func (writer *BufferWriter) Reset() {
	writer.line.Start = 0
	writer.line.End = 0
}

//删除分割字符
func (writer *BufferWriter) RemoveCutChar() {
	decrCount := 0
	size := len(writer.buffer)
	for i := size - 1; i >= 0; i-- {
		if writer.buffer[i] == writer.cutChar {
			decrCount++
			for j := i + 1; j < size; j++ {
				writer.buffer[j-1] = writer.buffer[j]
			}
		}
	}
	writer.buffer = writer.buffer[0 : size-decrCount]
}

//缓存容量
func (writer *BufferWriter) Cap() int {
	return cap(writer.buffer)
}

//缓存大小
func (writer *BufferWriter) Size() int {
	return len(writer.buffer)
}

func (writer *BufferWriter) SetCutChar(cutChar byte) {
	writer.cutChar = cutChar
}

func (writer *BufferWriter) ReplaceByte(old, new byte) {
	for i := len(writer.buffer) - 1; i >= 0; i-- {
		if writer.buffer[i] == old {
			writer.buffer[i] = new
		}
	}
}

func (writer *BufferWriter) Remove(start, length int) bool {
	if start < 0 || length <= 0 || start + length > writer.Size() {
		return false
	}
	size := writer.Size() - length
	for index := start; index < size; index++{
		writer.buffer[index] = writer.buffer[index + length]
	}
	writer.buffer = writer.buffer[:size]
	return true
}

func (writer *BufferWriter) RemoveByte(index int) bool{
	return writer.Remove(index, 1)
}

func (writer *BufferWriter) Recycle() {
	GetBufferPool().Add(writer)
}

