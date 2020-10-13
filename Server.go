package main

import (
	"fmt"

	"github.com/HyanSource/Helge/hinterface"
	"github.com/HyanSource/Helge/hnet"
	"github.com/HyanSource/Helge_Slot/slot"
	"github.com/golang/protobuf/proto"
)

type Play struct {
	hnet.BaseRouter
}

func (t *Play) Handle(request hinterface.IRequest) {
	result := slot.PlayGame.Result(slot.PlayGame.GetTable())
	r, err := proto.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.GetConnection().SendMsg(100, r)
}

func main() {
	s := hnet.NewServer()

	s.AddRouter(100, &Play{}) //遊玩業務

	s.Serve()
}
