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
	//nats原始消息
	GetRawMsg() *nats.Msg
	//序列化后的消息
	Message() Message
	//应当消息
	RespondMessage(msg interface{}) error
}

//actor消息传递上下文
type actorContext struct {
	rawMsg *nats.Msg
	msg interface{}
}

func newContext(raw *nats.Msg,msg Message) *actorContext {
	return &actorContext{raw,msg}
}

//nats消息
func (c *actorContext) GetRawMsg() *nats.Msg  {
	return c.rawMsg
}

//序列化后的消息
func (c *actorContext) Message() Message  {
	return c.msg
}

//应当消息
func  (c *actorContext) RespondMessage(msg interface{}) error  {
	data,err := encode.Encode(msg)
	if err != nil {
	    return err
	}
	return c.rawMsg.Respond(data)
}