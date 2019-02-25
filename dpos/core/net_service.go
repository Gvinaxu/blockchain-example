package core

import (
	"fmt"
	"math/rand"
	"time"
)

func Run() {
	delegateIndex = 0 // 初始委托者索引

	// 创世区块
	genesisBlock := Block{}
	genesisBlock = Block{0, time.Now().Unix(), "", CalculateBlockHash(genesisBlock), "", ""}
	blockChain = append(blockChain, genesisBlock)
	delegateIndex++

	countDelegate := len(delegates)

	for delegateIndex < countDelegate {
		time.Sleep(time.Second * 3)
		fmt.Println(delegateIndex)

		// 创建新的区块
		rand.Seed(int64(time.Now().Unix())) // 随机
		v := rand.Intn(100)
		end := blockChain[len(blockChain)-1]
		new := GenerateBlock(end, string(v), delegates[delegateIndex])

		fmt.Printf("Blockchain...%v\n", new)
		if IsBlockValid(new, end) {
			blockChain = append(blockChain, new)
		}

		delegateIndex = (delegateIndex + 1) % countDelegate
		if delegateIndex == 0 {
			delegates = RandDelegates(delegates) // 下一轮投票
		}
	}
}
