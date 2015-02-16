package np

import (
//	"errors"
//"fmt"
)

type StreamHandler struct {
	converter Converter

	struBuf []byte

	destruBuf []byte
	offset    int

	hasSlice bool

	endFlag []byte
}

func NewStreamHandler(endFlag []byte, converter Converter) DataHandler {
	sh := new(StreamHandler)
	sh.converter = converter
	sh.struBuf = make([]byte, 0)
	sh.offset = 0
	sh.hasSlice = false
	sh.endFlag = endFlag
	return sh
}

func (sh *StreamHandler) Load(data []byte) bool {
	sh.struBuf = append(sh.struBuf, data...)
	if compareSlice(sh.struBuf[len(sh.struBuf)-len(sh.endFlag):], sh.endFlag) {
		return true
	} else {
		return false
	}
}

func (sh *StreamHandler) ToObject() (interface{}, error) {
	temp := sh.struBuf[:len(sh.struBuf)-len(sh.endFlag)]
	object, err := sh.converter.ToObject(&temp)
	if err != nil {
		return nil, err
	}
	sh.struBuf = make([]byte, 0)
	return object, nil
}

func (sh *StreamHandler) Cut(object *interface{}) error {
	temp, err := sh.converter.ToArrayByte(object)
	if err != nil {
		return err
	}
	sh.destruBuf = temp
	sh.destruBuf = append(sh.destruBuf, sh.endFlag...)
	sh.offset = 0
	sh.hasSlice = true
	return nil
}

func (sh *StreamHandler) GetOneSlice() ([]byte, bool) {
	if len(sh.destruBuf) == sh.offset {
		return nil, false
	}

	if len(sh.destruBuf)-sh.offset > MaxBodySize {
		buf := sh.destruBuf[sh.offset : sh.offset+MaxBodySize]
		sh.offset += MaxBodySize
		return buf, true
	} else {
		buf := sh.destruBuf[sh.offset:]
		sh.offset = len(sh.destruBuf)
		sh.hasSlice = false
		return buf, false
	}
}

func (sh *StreamHandler) HasSlice() bool {
	return sh.hasSlice
}

func compareSlice(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

//func (sh *StreamHandler) Reset() {
//	sh.struBuf = make([]byte, 0)
//	sh.offset = 0
//	sh.hasSlice = false
//}
