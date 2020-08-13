package actor

import "github.com/magicsea/nactor/encode"

func init() {
	encode.RegisterName((*Kill)(nil))
	encode.RegisterName((*Watch)(nil))
	encode.RegisterName((*Unwatch)(nil))
	encode.RegisterName((*WatchTerminated)(nil))
	encode.RegisterName((*HeartBeat)(nil))
	encode.RegisterName((*Started)(nil))
}

//开始运行
type Started struct {
}

//杀掉actor
type Kill struct {
	Reason string
	Who string
}

//监听一个actor
type Watch struct {
	Watcher string
}

//解除监听一个actor
type Unwatch struct {
	Watcher string
}

//监听目标销毁
type WatchTerminated struct {
	Who string
}


//TODO:监听心跳
type HeartBeat struct {

}