package service

type MsgPack struct {
	MsgId   uint32
	DataLen uint32
	Data    []byte
}

func NewMsgPack(msgid uint32, data []byte) *MsgPack { //data has been compressed
	return &MsgPack{
		MsgId:   msgid,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
