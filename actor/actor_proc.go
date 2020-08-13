package actor


//actor上层处理器
type ActorProc interface {
	Receive(ctx Context)
	OnDestroy()
}
