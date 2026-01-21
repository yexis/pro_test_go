package engine

type BlockType int

const (
	// BlockTypeBlock 并行组等待
	BlockTypeBlock BlockType = iota
	// BlockTypeFinalBlock 全局等待
	BlockTypeFinalBlock
	// BlockTypeNoBlock 不等待
	BlockTypeNoBlock
)

// BlockJob ...
// used within concurrent-job
type BlockJob struct {
	ActionI
	needWait BlockType
}

func newBlockAction(a ActionI, nw BlockType) *BlockJob {
	return &BlockJob{ActionI: a, needWait: nw}
}

func (bj *BlockJob) NeedWait() BlockType {
	return bj.needWait
}
