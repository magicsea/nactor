package encode

//内置消息类型

const (
	Int8   = 2
	UInt8  = 3
	Int16  = 4
	UInt16 = 5
	Int32  = 6
	UInt32 = 7
	Int64  = 8
	UInt64 = 9
	String = 10
	Bytes  = 11

	Bool           = 14
	Float32        = 15
	Float64        = 16

	//protobuff结构
	Proto = 50
	//gob结构
	Gob = 51
)
