package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

var (
	difficulty = int64(1)
)

func GenerateBlock(old Block, value string) Block {
	var new Block
	new.Height = old.Height + 1
	new.Timestamp = time.Now().Unix()
	new.Difficulty = difficulty
	new.Value = value
	new.PreHash = old.Hash

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		new.Nonce = hex
		newHash := calculateHash(new)
		if !isHashValid(newHash, int(difficulty)) {
			time.Sleep(time.Millisecond)
		} else {
			difficulty += 2
			new.Hash = newHash
			break
		}
	}
	return new
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	hasPrefix := strings.HasPrefix(hash, prefix)
	return hasPrefix
}

func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%d%s%s%s", block.Height, block.Timestamp, block.Value, block.PreHash, block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isBlockValid(new, old Block) bool {
	if old.Height+1 != new.Height {
		return false
	}
	if old.Hash != new.PreHash {
		return false
	}
	if calculateHash(new) != new.Hash {
		return false
	}
	return true
}
