package main

import (
	"fmt"

	"github.com/HyanSource/Helge_Slot/metrics"

	"github.com/HyanSource/Helge/hinterface"
	"github.com/HyanSource/Helge/hnet"
	"github.com/HyanSource/Helge_Slot/pb"
	"github.com/HyanSource/Helge_Slot/player"
	"github.com/HyanSource/Helge_Slot/slot"
	"github.com/golang/protobuf/proto"
)

var (
	Metrics *metrics.LocalMetrics
)

type Play struct {
	hnet.BaseRouter
}

func (t *Play) Handle(request hinterface.IRequest) {

	play := &pb.Play{}
	err := proto.Unmarshal(request.GetMessage().GetData(), play)
	if err != nil {
		fmt.Println(err)
		return
	}
	p, err := request.GetConnection().GetPropertys().GetProperty("player")
	if err != nil {
		fmt.Println(err)
		return
	}

	result := slot.PlayGame.Result(slot.PlayGame.GetTable())

	money := int32(p.(*player.Player).Money) - play.Bet + play.Bet*result.Odds
	p.(*player.Player).Money = int(money)
	result.Money = &pb.Money{Money: money}

	{
		Metrics.AllPlayCounter.Inc(1)
		Metrics.WinMoneyCounter.Inc(1)
	}

	r, err := proto.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.GetConnection().SendMsg(100, r)
}

func StartConnection(conn hinterface.IConnection) {
	const InitMoney = 10000
	money := &pb.Money{Money: InitMoney}
	m, err := proto.Marshal(money)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.SendMsg(10, m)
	conn.GetPropertys().SetProperty("player", &player.Player{Money: InitMoney})
}

func main() {

	Metrics = metrics.NewLocalMetrics()

	s := hnet.NewServer()
	s.GetHook().SetHook("start", StartConnection)
	s.AddRouter(100, &Play{}) //遊玩業務

	s.Serve()
}
