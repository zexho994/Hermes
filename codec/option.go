package codec

const MagicNumber = 0x994994

// Option 协商时候使用的
type Option struct {
	// 表示为Hermes RPC 请求
	MagicNumber int
	// 编码方式
	CodeType Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodeType:    GobType,
}
