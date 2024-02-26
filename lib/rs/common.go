package rs

const (
	DataShard     = 4
	ParityShard   = 2
	AllSharad     = DataShard + ParityShard
	BlockPerShard = 8000
	BlockSie      = BlockPerShard * DataShard
)
