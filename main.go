package main

import (
	"fmt"

	bw "gopkg.in/immesys/bw2bind.v5"
)

func main() {
	fmt.Printf("connecting to BW\n")
	bwc := bw.ConnectOrExit("")
	fmt.Printf("setting entity\n")
	bwc.SetEntityFromEnvironOrExit()
	fmt.Printf("connected and set entity ok\n")
	RealURI := "amplab/sensors/s.hamilton/00126d0700000061/i.temperature/signal/operative"
	FakeURI := "amplab/sensors/s.hamilton/00126d0700000065/i.temperature/signal/operative"
	_ = RealURI
	_ = FakeURI
	subchan := bwc.SubscribeOrExit(&bw.SubscribeParams{
		URI:       RealURI,
		AutoChain: true,
	})
	fmt.Printf("subscribed to hamiltons ok\n")
	for m := range subchan {
		go procMsg(bwc, m)
	}
	fmt.Printf("subscription channel ended\n")
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
	fmt.Printf("got a message\n")
	ham := m.GetOnePODF("2.0.11.2")
	if ham == nil {
		fmt.Printf("No PO found, skipping\n")
		return
	}
	hamdata := HamiltonData{}
	ham.(bw.MsgPackPayloadObject).ValueInto(&hamdata)
	fmt.Printf("hamilton data is %#v\n", hamdata)
	if hamdata == LastData {
		fmt.Printf("data is the same, skipping\n")
		return
	}
	LastData = hamdata
	//Feel free to change this, but this will kinda work
	lc := &LifxCommand{
		Hue:   4.0 * (float64(hamdata.Buttons%4) / 4.0),
		Sat:   1.0,
		Bri:   1.0,
		State: true,
	}
	po, err := bw.CreateMsgPackPayloadObject(bw.PONumHSBLightMessage, lc)
	fmt.Printf("create LIFX PO err=%v\n", err)
	bwc.PublishOrExit(&bw.PublishParams{
		URI:            "ucberkeley/eop/lifx/s.lifx/0/i.hsb-light",
		PayloadObjects: []bw.PayloadObject{po},
		AutoChain:      true,
	})
}
