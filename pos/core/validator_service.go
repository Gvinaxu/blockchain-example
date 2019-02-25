package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func HandleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	var address string
	io.WriteString(conn, "Enter token balance:")

	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number:%v", scanBalance.Text(), err)
			return
		}
		address = calculateHash(time.Now().String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	io.WriteString(conn, "\nEnter a new value:")
	scanBPM := bufio.NewScanner(conn)
	go func() {
		for {
			for scanBPM.Scan() {
				value := scanBPM.Text()

				mutex.Lock()
				end := blockChain[len(blockChain)-1]
				mutex.Unlock()
				new := generateBlock(end, value, address)

				if isBlockValid(new, end) {
					candidateBlocks <- new
				}
				io.WriteString(conn, "\nEnter a new value:")
			}
		}
	}()

	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output, err := json.Marshal(blockChain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}
}

// 选择可以出块的生产者
func PickWinner() {
	time.Sleep(time.Second * 14)
	mutex.Lock()
	temp := tempBlocks //等待被选中
	mutex.Unlock()

	var rewardPool []string // 奖池 账户拥有n个token，就在其中记录n次该地址
	if len(temp) > 0 {
	OUT:
		for _, block := range temp {
			for _, node := range rewardPool {
				if block.Validator == node {
					continue OUT
				}
			}

			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			if k, ok := setValidators[block.Validator]; ok {
				// k 代表的是当前验证者的tokens
				for i := 0; i < k; i++ {
					rewardPool = append(rewardPool, block.Validator)
				}
				fmt.Println(rewardPool)
			}
		}

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		winner := rewardPool[r.Intn(len(rewardPool))]

		for _, block := range temp {
			if block.Validator == winner {
				mutex.Lock()
				blockChain = append(blockChain, block)
				mutex.Unlock()
				for range validators {
					announcements <- "\nwinning validator: " + winner + "\n"
				}
				break
			}
		}
	}
	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}
