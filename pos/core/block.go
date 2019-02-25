package core

import "sync"

type Block struct {
	Height    int64  `json:"height"`
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
	Validator string `json:"validator"`
}

var (
	blockChain      []Block
	tempBlocks      []Block
	candidateBlocks = make(chan Block)  // 提议区块
	announcements   = make(chan string) // 广播信息
	mutex           = &sync.Mutex{}
	validators      = make(map[string]int)
)
