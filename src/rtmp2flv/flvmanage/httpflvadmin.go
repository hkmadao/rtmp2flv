package flvmanage

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/httpflvmanage"
)

var hfas *HttpflvAdmin

type HttpflvAdmin struct {
	hfms sync.Map
}

func init() {
	hfas = &HttpflvAdmin{}
}

func GetSingleHttpflvAdmin() *HttpflvAdmin {
	return hfas
}

// func (hfa *HttpflvAdmin) ExistsHttpFlvManager(code string) bool {
// 	_, b := hfa.hfms.Load(code)
// 	return b
// }

func (hfa *HttpflvAdmin) AddHttpFlvManager(
	pktStream <-chan av.Packet,
	code string,
	codecs []av.CodecData,
) {
	hfm := httpflvmanage.NewHttpFlvManager(pktStream, code, codecs)
	hfa.hfms.Store(code, hfm)
}

func (hfa *HttpflvAdmin) StopWrite(code string) {
	v, ok := hfa.hfms.Load(code)
	if ok {
		ffw := v.(*httpflvmanage.HttpFlvManager)
		ffw.StopWrite()
	}
}

func (hfa *HttpflvAdmin) StartWrite(code string) {
	v, ok := hfa.hfms.Load(code)
	if ok {
		ffw := v.(*httpflvmanage.HttpFlvManager)
		ffw.StopWrite()
		hfa.AddHttpFlvManager(ffw.GetPktStream(), code, ffw.GetCodecs())
	}
}

//添加播放者
func (hfa *HttpflvAdmin) AddHttpFlvPlayer(
	playerDone <-chan interface{},
	pulseInterval time.Duration,
	code string,
	writer io.Writer,
) (<-chan interface{}, error) {
	v, b := hfa.hfms.Load(code)
	if b {
		hfm := v.(*httpflvmanage.HttpFlvManager)
		return hfm.AddHttpFlvPlayer(playerDone, pulseInterval, writer)
	}
	return nil, errors.New("camera no connection")
}
