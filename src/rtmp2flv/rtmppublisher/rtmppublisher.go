package rtmppublisher

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/flvmanage"
)

type RtmpServer interface {
	Load(key interface{}) (interface{}, bool)
	Store(key, value interface{})
	Delete(key interface{})
}

type Publisher struct {
	code         string
	codecs       []av.CodecData
	connDone     <-chan interface{}
	pktStream    <-chan av.Packet
	ffmPktStream <-chan av.Packet
	hfmPktStream <-chan av.Packet
	rtmpserver   RtmpServer
}

func NewPublisher(connDone <-chan interface{}, pktStream <-chan av.Packet, code string, codecs []av.CodecData, rs RtmpServer) *Publisher {
	r := &Publisher{
		connDone:     connDone,
		pktStream:    pktStream,
		code:         code,
		codecs:       codecs,
		ffmPktStream: make(chan av.Packet, 1024),
		hfmPktStream: make(chan av.Packet, 1024),
		rtmpserver:   rs,
	}
	r.pktTransfer()
	return r
}

func (r *Publisher) Done() {
	<-r.connDone
}

func (r *Publisher) pktTransfer() {
	ffmPktStream, hfmPktStream := tee(r.connDone, r.pktStream)
	r.ffmPktStream = ffmPktStream
	r.hfmPktStream = hfmPktStream
	logs.Debug("publisher [%s] create customer", r.code)
	flvmanage.GetSingleFileFlvManager().FlvWrite(r.ffmPktStream, r.code, r.codecs)
	flvmanage.GetSingleHttpflvAdmin().AddHttpFlvManager(r.hfmPktStream, r.code, r.codecs)
}

func tee(done <-chan interface{}, in <-chan av.Packet) (<-chan av.Packet, <-chan av.Packet) {
	//设置缓冲，调节前后速率
	out1 := make(chan av.Packet, 1024)
	out2 := make(chan av.Packet, 1024)
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			var out1, out2 = out1, out2 // 私有变量覆盖
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil // 置空阻塞机制完成select轮询
				case out2 <- val:
					out2 = nil
				default:
					logs.Error("publisher tee lose packet")
				}
			}
		}
	}()
	return out1, out2
}
