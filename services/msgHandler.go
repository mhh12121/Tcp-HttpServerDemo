package service

type MsgHandler struct {
	Handler        map[uint32]func()
	WorkerPoolSize uint32
	// TaskQueue []chan

}

func NewMsgHandler() *MsgHandler { //data has been compressed
	return &MsgHandler{}
}
