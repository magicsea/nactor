package actor


//actor上层处理器
type ActorProc interface {
	OnStart()
	Receive(ctx Context)
	OnDestroy()
}
