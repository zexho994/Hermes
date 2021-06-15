package codec

import (
	"io"
)

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(closer io.ReadWriteCloser) Codec

// Header stores data other than request parameters and response data.
type Header struct {
	// method's name. format like "<service>.<method>"
	ServiceMethod string
	// request-id, the purpose is to distinguish the between requests
	Seq uint64
	// Save the  server error
	Error string
}

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
	//NewCodecFuncMap[JsonType] = NewGobCodec
}
