package main


type IteratorLine interface {
	MoveNext() bool
	Current() string
	Reset()
}
