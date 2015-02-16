package np

import (
	"crypto/md5"
	//"errors"
)

type ChanidConverter struct {
}

func (cc *ChanidConverter) ToObject(data *[]byte) (interface{}, error) {
	if len(*data) != md5.Size {
		return nil, DataConvertErr
	}
	return [md5.Size]byte{
		(*data)[0], (*data)[1], (*data)[2], (*data)[3],
		(*data)[4], (*data)[5], (*data)[6], (*data)[7],
		(*data)[8], (*data)[9], (*data)[10], (*data)[11],
		(*data)[12], (*data)[13], (*data)[14], (*data)[15]}, nil
}

func (cc *ChanidConverter) ToArrayByte(object *interface{}) ([]byte, error) {
	array, ok := (*object).([md5.Size]byte)
	if !ok {
		return nil, DataConvertErr
	}
	return array[:], nil
}
