package core

type Block struct {
	Height    int64  `json:"height"`
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
	Hash      string `json:"hash"`
	PreHash   string `json:"pre_hash"`
	Delegate  string `json:"delegate"`
}

var (
	blockChain    []Block
	delegates     = []string{"001", "002", "003", "004", "005", "006", "007"}
	delegateIndex int
)
