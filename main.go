package main

import (
	"fmt"

	bw "gopkg.in/immesys/bw2bind.v5"
)

func main() {
	bwc := bw.ConnectOrExit("")
	bwc.SetEntityFromEnvironOrExit()
	subchan := bwc.SubscribeOrExit(&bw.SubscribeParams{
		URI:       "amplab/sensors/s.hamilton/00126d0700000061/i.temperature/signal/operative",
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

func procMsg(bwc *bw.BW2Client, m *bw.SimpleMessage) {
	ham := m.GetOnePODF("2.0.11.2")
	if ham == nil {
		return
	}
	hamdata := HamiltonData{}
	ham.(bw.MsgPackPayloadObject).ValueInto(&hamdata)
	fmt.Printf("hamilton data is %v\n", hamdata)
	lc := &LifxCommand{
		Hue:   0.0, /*JACK DO THIS*/
		Sat:   0.0, /*JACK DO THIS*/
		Bri:   0.0, /*JACK DO THIS*/
		State: true,
	}
	po, _ := bw.CreateMsgPackPayloadObject(bw.PONumHSBLightMessage, lc)
	bwc.PublishOrExit(&bw.PublishParams{
		URI:            "ucberkeley/eop/lifx/s.lifx/0/i.hsb-light",
		PayloadObjects: []bw.PayloadObject{po},
	})
}
