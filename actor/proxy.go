package actor

import (
	"github.com/magicsea/nactor/encode"
	"github.com/nats-io/nats.go"
	"time"
)

/**
 *actor 代理,
 *name可以是actor的名字，也可以是actor额外订阅的频道
 *
 *
**/
type Proxy struct {
	name string
	conn *nats.Conn
}

//NewProxy
func NewProxy(name string,conn *nats.Conn) *Proxy  {
	return &Proxy{name:name,conn:conn}
}

//通知消息
func (p *Proxy) Tell(msg Message) error {
	data,err := encode.Encode(msg)
	if err != nil {
	    return err
	}
	return p.conn.Publish(p.name,data)
}

//同步请求应答
func (p *Proxy) Request(msg Message,timeout time.Duration) (Message,error) {
	data,err := encode.Encode(msg)
	if err != nil {
		return nil,err
	}
	rsp,errRsp := p.conn.Request(p.name,data,timeout)
	if errRsp != nil {
	    return nil,errRsp
	}
	rmsg,errDe := encode.Decode(rsp.Data)
	if errDe != nil {
	    return nil,errDe
	}
	return rmsg,nil
}


//异步请求应答
type AsyncRequestRsp struct {
	RawMsg *nats.Msg
	Msg Message
	Err error
}
/**
 *异步请求应答
 *返回:处理结果的chan
 *
**/
func (p *Proxy) AsyncRequest(msg Message,timeout time.Duration) (chan AsyncRequestRsp,error) {
	data,err := encode.Encode(msg)
	if err != nil {
		return nil,err
	}
	ch := make(chan AsyncRequestRsp)
	go func(ch chan AsyncRequestRsp) {
		m,e := p.conn.Request(p.name,data,timeout)
		if e != nil {
			ch<-AsyncRequestRsp{nil,nil,e}
			return
		}
		data := m.Data
		rmsg,errDe := encode.Decode(data)
		if errDe != nil {
			ch<-AsyncRequestRsp{nil,nil,errDe}
			return
		}
		ch<-AsyncRequestRsp{m,rmsg,nil}
	}(ch)

	return ch,nil
}
