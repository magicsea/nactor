package actor

import (
	"github.com/magicsea/nactor/encode"
	"github.com/nats-io/nats.go"
)

//消息类型
type Message = interface{}

/**
 *actor消息传递上下文接口
 *隐藏actorContext内部实现
 *
**/
type Context interface {
	//自身actor名字
	ActorName() string
	//nats原始消息
	GetRawMsg() *nats.Msg
	//序列化后的消息
	Message() Message

	//请求发送者
	RequestSender() string
	//应答消息
	RespondMessage(msg interface{}) error
	//监听
	Watch(target string)
	//解除监听
	Unwatch(target string)
	//流转消息
	Forward(target string)
}

//actor消息传递上下文
type actorContext struct {
	actorName string
	conn *nats.Conn
	rawMsg *nats.Msg
	msg interface{}

}

func newContext(actor *Actor,raw *nats.Msg,msg Message) *actorContext {
	return &actorContext{
		actorName:actor.name,
		rawMsg:raw,
		msg:msg,
		conn:actor.conn,
	}
}

//name
func (c *actorContext) ActorName() string  {
	return c.actorName
}

//nats消息
func (c *actorContext) GetRawMsg() *nats.Msg  {
	return c.rawMsg
}

//序列化后的消息
func (c *actorContext) Message() Message  {
	return c.msg
}

//请求发送者
func (c *actorContext) RequestSender() string {
	return c.rawMsg.Reply
}

//应答消息
func  (c *actorContext) RespondMessage(msg interface{}) error  {
	data,err := encode.Encode(msg)
	if err != nil {
	    return err
	}

	return c.rawMsg.Respond(data)
}

//监听
func  (c *actorContext) Watch(target string) {
	p := NewProxy(target,c.conn)
	p.Tell(&Watch{Watcher:c.ActorName()})
}
//解除监听
func  (c *actorContext) Unwatch(target string) {
	p := NewProxy(target,c.conn)
	p.Tell(&Unwatch{Watcher:c.ActorName()})
}
//流转消息
func  (c *actorContext) Forward(target string) {
	p := NewProxy(target,c.conn)
	p.SendRaw(c.rawMsg.Data)
}