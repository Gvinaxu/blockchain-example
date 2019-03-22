package core

import "sync"

func init() {
	block := Block{Height: 1, Hash: "ebd8266b1d2f357234443451c62e2244cf154be4471dc84a634f080ec9bf178f"}
	blockChain = append(blockChain, block)
}

type Block struct {
	Height     int64  `json:"height"`
	Timestamp  int64  `json:"timestamp"`
	Difficulty float64  `json:"difficulty"`
	Value      string `json:"value"`
	Hash       string `json:"hash"`
	PreHash    string `json:"pre_hash"`
	Nonce      string `json:"nonce"`
}

var blockChain []Block

type Message struct {
	Value string `json:"value"`
}

var Mutex = &sync.Mutex{}
