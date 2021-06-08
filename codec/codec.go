package codec

type Header struct {
	// 方法名称
	serviceMethod string
	// 请求序号,可以理解为请求的ID
	seq uint64
	// 错误信息
	error string
}
