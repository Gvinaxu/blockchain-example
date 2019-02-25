package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func generateBlock(old Block, value string, address string) Block {
	var new Block

	new.Height = old.Height + 1
	new.Timestamp = time.Now().Unix()
	new.Value = value
	new.PrevHash = old.Hash
	new.Hash = CalculateBlockHash(new)
	new.Validator = address
	return new
}

func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CalculateBlockHash(block Block) string {
	record := fmt.Sprintf("%d%d%s%s", block.Height, block.Timestamp, block.Value, block.PrevHash)
	return calculateHash(record)
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Height+1 != newBlock.Height {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if CalculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true

}
