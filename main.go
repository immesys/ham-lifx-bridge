package main

import (
	"fmt"

	bw "gopkg.in/immesys/bw2bind.v5"
)

func main() {
	bwc := bw.ConnectOrExit("")
	bwc.SetEntityFromEnvironOrExit()
	RealURI := "amplab/sensors/s.hamilton/00126d0700000061/i.temperature/signal/operative"
	FakeURI := "amplab/sensors/s.hamilton/00126d0700000065/i.temperature/signal/operative"
	_ = RealURI
	_ = FakeURI
	subchan := bwc.SubscribeOrExit(&bw.SubscribeParams{
		URI:       FakeURI,
		AutoChain: true,
	})
	for m := range subchan {
		procMsg(bwc, m)
	}
}

type LifxCommand struct {
	Hue   float64 `msgpack:"hue"`
	Sat   float64 `msgpack:"sat"`
	Bri   float64 `msgpack:"bri"`
	State bool    `msgpack:"state"`
}
type HamiltonData struct {
	Lux     float64 `msgpack:"lux"`
	Buttons int     `msgpack:"button_events"`
	Temp    float64 `msgpack:"air_temp"`
	RH      float64 `msgpack:"air_rh"`
}

var LastData HamiltonData

func procMsg(bwc *bw.BW2Client, m *bw.SimpleMessage) {
	ham := m.GetOnePODF("2.0.11.2")
	if ham == nil {
		return
	}
	hamdata := HamiltonData{}
	ham.(bw.MsgPackPayloadObject).ValueInto(&hamdata)
	fmt.Printf("hamilton data is %#v\n", hamdata)
	if hamdata == LastData {
		return
	}
	LastData = hamdata
	//Feel free to change this, but this will kinda work
	lc := &LifxCommand{
		Hue:   4.0 * float64(hamdata.Buttons%4),
		Sat:   1.0,
		Bri:   1.0,
		State: true,
	}
	po, _ := bw.CreateMsgPackPayloadObject(bw.PONumHSBLightMessage, lc)
	bwc.PublishOrExit(&bw.PublishParams{
		URI:            "ucberkeley/eop/lifx/s.lifx/0/i.hsb-light",
		PayloadObjects: []bw.PayloadObject{po},
	})
}
