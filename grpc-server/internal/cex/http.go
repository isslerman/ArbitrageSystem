package cex

type HTTPMethod int

const (
	GET  HTTPMethod = iota // auto enum - first is 0
	POST                   // 1
	PUT                    // 2
	DELETE
	PATCH
	HEAD
	OPTIONS
	TRACE
	CONNECT
)
