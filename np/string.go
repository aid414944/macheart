package np

type StringConverter struct {
}

func (sc *StringConverter) ToObject(data *[]byte) (interface{}, error) {
	return string(*data), nil
}

func (sc *StringConverter) ToArrayByte(object *interface{}) ([]byte, error) {
	return []byte((*object).(string)), nil
}
