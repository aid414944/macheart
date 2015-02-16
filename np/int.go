package np

import (
	"encoding/binary"
	//"errors"
)

type IntConverter struct {
}

func (ic *IntConverter) ToObject(data *[]byte) (interface{}, error) {
	temp := binary.LittleEndian.Uint64(*data)
	return temp, nil
}

func (ic *IntConverter) ToArrayByte(object *interface{}) ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64((*object).(int)))
	return buf, nil
}
