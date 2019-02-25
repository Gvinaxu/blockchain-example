package core

import (
	"log"
	"net"
	"time"
)

func Run() {
	genesisBlock := Block{}
	genesisBlock = Block{Timestamp: time.Now().Unix(), Hash: CalculateBlockHash(genesisBlock)}
	// 创世区块
	blockChain = append(blockChain, genesisBlock)
	port := ":2345"
	server, err := net.Listen("tcp", port)
	defer server.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP Server Listening on port ", port)

	//
	go func() {
		for temp := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, temp)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			PickWinner()
		}
	}()

	for {
		conn, err := server.Accept() //tcp接收监听
		if err != nil {
			log.Fatal(err)
		}
		go HandleConn(conn) //
	}
}
