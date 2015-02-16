package np

import (
//"fmt"
//"np/conf"
)

type BlockHandler struct {
	converter Converter

	struBuf []byte

	destruBuf []byte
	offset    int

	hasSlice bool

	size uint
}

func NewBlockHandler(size uint, converter Converter) DataHandler {
	bh := new(BlockHandler)
	bh.converter = converter
	bh.struBuf = make([]byte, 0)
	bh.offset = 0
	bh.hasSlice = false
	bh.size = size
	return bh
}

func (bh *BlockHandler) Load(data []byte) bool {
	bh.struBuf = append(bh.struBuf, data...)
	if uint(len(bh.struBuf)) >= bh.size { // 大于的情况下是不正确的，说明另一端在发送无效数据，这时候应该关闭链接，现在暂时不进行处理，等进入转换阶段再处理
		return true
	} else {
		return false
	}
}

func (bh *BlockHandler) ToObject() (interface{}, error) {
	temp, err := bh.converter.ToObject(&(bh.struBuf))
	if err != nil {
		return nil, err
	}
	bh.struBuf = make([]byte, 0)
	return int(temp.(uint64)), nil
}

func (bh *BlockHandler) Cut(object *interface{}) error {
	//bh.destruBuf = make([]byte, bh.size)
	temp, err := bh.converter.ToArrayByte(object)
	if err != nil {
		return err
	}
	bh.destruBuf = temp
	bh.offset = 0
	bh.hasSlice = true
	return nil
}

func (bh *BlockHandler) GetOneSlice() ([]byte, bool) {
	if len(bh.destruBuf) == bh.offset {
		return nil, false
	}

	if len(bh.destruBuf)-bh.offset > MaxBodySize {
		buf := bh.destruBuf[bh.offset : bh.offset+MaxBodySize]
		bh.offset += MaxBodySize
		return buf, true
	} else {
		buf := bh.destruBuf[bh.offset:]
		bh.offset = len(bh.destruBuf)
		bh.hasSlice = false
		return buf, false
	}
}

func (bh *BlockHandler) HasSlice() bool {
	return bh.hasSlice
}
