package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateBlock(old Block, value, address string) Block {
	var new Block
	new.Height = old.Height + 1
	new.Timestamp = time.Now().Unix()
	new.Value = value
	new.PreHash = old.PreHash
	new.Hash = CalculateBlockHash(new)
	new.Delegate = address
	return new
}

func IsBlockValid(new, old Block) bool {
	if old.Height+1 != new.Height {
		return false
	}
	if old.Hash != new.PreHash {
		return false
	}
	if CalculateBlockHash(new) != new.Hash {
		return false
	}
	return true
}

func CalculateBlockHash(block Block) string {
	record := fmt.Sprintf("%d%d%s%s", block.Height, block.Timestamp, block.Value, block.PreHash)
	return calculateHash(record)
}

func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
