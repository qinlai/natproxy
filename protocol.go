package main

type ProtocolFormat struct {
}

func NewProtocolFormat() ProtocolFormat {
	return ProtocolFormat{}
}

func (ProtocolFormat) DecodeMessage(b []byte) (interface{}, error) {
	return b, nil
}

func (ProtocolFormat) EncodeMessage(msg interface{}) ([]byte, error) {
	return msg.([]byte), nil
}
