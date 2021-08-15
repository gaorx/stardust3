package sdsnowflake

var (
	DefaultNode *Node
	ZeroNode    *Node
)

func init() {
	defaultNode, err := NewFromIP()
	if err != nil {
		panic(err)
	}
	DefaultNode = defaultNode

	zeroNode, err := New(0)
	if err != nil {
		panic(err)
	}
	ZeroNode = zeroNode
}

func GenerateByIP() int64 {
	return DefaultNode.Generate()
}

func GenerateZero() int64 {
	return ZeroNode.Generate()
}
