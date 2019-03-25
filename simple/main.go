package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	httpPort     string
	p2pPort      string
	listPeers    = make(map[string]bool)
	BlockChain   = make([]Block, 0)
	orphanBlocks = make(map[OrphanBlockKey]Block)
	rwMutex      sync.RWMutex
	version      = "0.1"
)

func main() {
	initBlockChain()
	go initHttpServer()
	initP2PServer()
}

type Block struct {
	Version   string `json:"version"`
	Hash      string `json:"hash"`
	Height    uint64 `json:"height"`
	PrevHash  string `json:"previous_hash"`
	TimeStamp int64  `json:"time_stamp"`
	MetaData  string `json:"meta_data"`
}

type OrphanBlockKey struct {
	height uint64 `json:"height"`
	hash   string `json:"hash"`
}

type Message struct {
	From     string            `json:"from"`
	To       string            `json:"to"`
	Type     string            `json:"type"`
	Msg      map[string][]byte `json:"msg"`
	MyHeight uint64            `json:"my_height"`
}

type Peer struct {
	Address string `json:"address"`
}

type BlockMetaData struct {
	MetaData string `json:"meta_data"`
}

func genesisBlock() (block Block) {
	block = Block{version, "", 0, "", time.Now().Unix(), "first block"}
	block.Hash = cryptoHash(&block)
	return
}

func initBlockChain() {
	BlockChain = append(BlockChain, genesisBlock())
}

func initHttpServer() {
	http.HandleFunc("/addPeer", addPeer)
	http.HandleFunc("/getPeers", getPeers)
	http.HandleFunc("/mineNewBlock", mineNewBlock)
	http.HandleFunc("/getBlockChain", getBlockChain)

	httpPort = "9090"
	if http_port := os.Getenv("HTTP_PORT"); http_port != "" {
		httpPort = http_port
	}
	log.Printf("Http server started on :%s", httpPort)

	var addr string
	addr = ":" + httpPort
	log.Fatal(http.ListenAndServe(addr, nil))
}

func initP2PServer() {
	p2pPort = "7676"
	if p2p_port := os.Getenv("P2P_PORT"); p2p_port != "" {
		p2pPort = p2p_port
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", getMyAddress())
	if err != nil {
		log.Println(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
	}
	log.Printf("P2PServer started on :%s", p2pPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go receive(conn)
	}
}

func addPeer(w http.ResponseWriter, r *http.Request) {
	peer := new(Peer)
	if r.Body == nil {
		http.Error(w, "Please send a request with body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(peer)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("A client calls http server to add a peer with address: %s", peer.Address)

	ip, err := net.ResolveTCPAddr("tcp4", peer.Address)
	if err != nil {
		log.Printf("Reveived wrong tcp address: %s", peer.Address)
		w.Write([]byte("Wrong tcp address!"))
		return
	}

	listPeers[ip.String()] = true
	log.Printf("Updated my peer list: %v\n", listPeers)

	msg := new(Message)
	msg.From = getMyAddress()
	msg.To = peer.Address
	msg.Type = "HELLO"
	msg.MyHeight = uint64(len(BlockChain))

	send(msg)
}

func getPeers(w http.ResponseWriter, r *http.Request) {
	listpeers_bytes, err := json.MarshalIndent(listPeers, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.Write(listpeers_bytes)
	w.Write([]byte("\n"))
}

func mineNewBlock(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request with body", 400)
		return
	}
	metaData := new(BlockMetaData)
	err := json.NewDecoder(r.Body).Decode(metaData)
	if err != nil {
		log.Println(err)
		return
	}

	rwMutex.Lock()
	block := new(Block)
	block.Version = version
	last := BlockChain[len(BlockChain)-1]
	block.PrevHash = last.Hash
	block.TimeStamp = time.Now().Unix()
	block.MetaData = metaData.MetaData
	hash := cryptoHash(block)
	block.Hash = hash
	block.Height = last.Height + 1
	BlockChain = append(BlockChain, *block)
	rwMutex.Unlock()

	broadcastBlock(block)
}

func getBlockChain(w http.ResponseWriter, r *http.Request) {

	data, err := json.MarshalIndent(BlockChain, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.Write(data)
	w.Write([]byte("\n"))
}

func getMyAddress() (addr string) {
	name, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	address, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		log.Println(err)
	}

	addr = address.String() + ":" + p2pPort
	return
}

func createConnection(addr string) (conn net.Conn) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Println(err)
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err)
		return nil
	}
	return conn
}

func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func cryptoHash(block *Block) (hash string) {
	record := fmt.Sprintf("%d%d%s%s", block.Height, block.TimeStamp, block.MetaData, block.PrevHash)
	return calculateHash(record)
}

func broadcastBlock(block *Block) {
	msg := new(Message)
	msg.From = getMyAddress()
	msg.Type = "BROADCAST_BLOCK"
	msg.MyHeight = uint64(len(BlockChain))
	msg.Msg = make(map[string][]byte)
	data, err := json.Marshal(block)
	if err != nil {
		log.Println(err)
	}
	msg.Msg["broadcast_block"] = data

	for to, _ := range listPeers {
		msg.To = to
		send(msg)
	}
}

func receive(conn net.Conn) {
	defer conn.Close()

	dec := json.NewDecoder(conn)
	msg := new(Message)
	err := dec.Decode(msg)
	if err != nil {
		log.Printf("Decoding error: %v", err)
		return
	}
	log.Printf("Received a p2p message: \n%v", *msg)

	handleMessage(msg)
}

func send(msg *Message) {
	conn := createConnection(msg.To)
	if conn == nil {
		return
	}
	enc := json.NewEncoder(conn)
	enc.Encode(msg)
	log.Printf("Sent a p2p message: \n%v", *msg)
}

func handleMessage(msg *Message) {
	switch msg.Type {
	case "HELLO":
		handleHello(msg)
	case "RESPONSE_HELLO":
		if msg.MyHeight > uint64(len(BlockChain)) {
			syncGetBlocks()
		}
	case "BROADCAST_BLOCK":
		handleBroadcastBlock(msg)
	case "BROADCAST_UPDATED":
		if msg.MyHeight > uint64(len(BlockChain)) {
			syncGetBlocks()
		}
	case "SYNC_GET_BLOCKS":
		handleSyncGetBlocks(msg)
	case "SYNC_BLOCKS":
		handleSyncBlocks(msg)
	}
}

func handleHello(msg *Message) {
	listPeers[msg.From] = true
	log.Printf("Updated my peers list after receiving a hello message: \n%v", listPeers)

	msg.To = msg.From
	msg.From = getMyAddress()
	msg.Type = "RESPONSE_HELLO"

	if msg.MyHeight > uint64(len(BlockChain)) { // A sender has a higher block height, which means a longer BlockChain.
		go syncGetBlocks()
	}

	msg.MyHeight = uint64(len(BlockChain))
	send(msg)
}

func handleBroadcastBlock(msg *Message) {
	defer rwMutex.Unlock()
	rwMutex.Lock()
	if msg.MyHeight <= uint64(len(BlockChain)) { // To avoid endless broadcast loop, refuse lagged block as well.
		return
	}

	broadcast_block := new(Block)
	json.Unmarshal(msg.Msg["broadcast_block"], broadcast_block)

	if msg.MyHeight > uint64(len(BlockChain)+1) { // Can't tell a new block's validity, then just store it into orphan block map.
		orphan_block_key := OrphanBlockKey{msg.MyHeight, cryptoHash(broadcast_block)}
		orphanBlocks[orphan_block_key] = *broadcast_block
		go syncGetBlocks()
		return
	}
	last := BlockChain[len(BlockChain)-1]
	if broadcast_block.PrevHash == last.Hash { //The new block is validity.
		BlockChain = append(BlockChain, *broadcast_block)
		log.Printf("Updated BlockChain after receiving a new broadcast_blcok, result BlockChain: \n%v", BlockChain)
		broadcastBlock(broadcast_block)
	}
}

func handleSyncGetBlocks(msg *Message) {
	defer rwMutex.RUnlock()
	rwMutex.RLock()

	if msg.MyHeight < uint64(len(BlockChain)) {
		hash := BlockChain[len(BlockChain)-1].Hash
		my_bytes, err := json.Marshal(hash)
		if err != nil {
			log.Println(err)
		}
		msg.Msg["hash_of_latest_block"] = my_bytes

		height := uint64(len(BlockChain))
		count := height - msg.MyHeight
		blocks := make(map[string]Block)
		for i := count; i > 0; i = i - 1 {
			for _, v := range BlockChain {
				if v.Hash == hash {
					blocks[hash] = v
					hash = v.PrevHash
				}
			}
		}
		my_bytes, err = json.Marshal(blocks)
		if err != nil {
			log.Println(err)
		}
		msg.Msg["sync_blocks"] = my_bytes

		msg.To = msg.From
		msg.From = getMyAddress()
		msg.Type = "SYNC_BLOCKS"
		msg.MyHeight = height

		send(msg)
	}
}

func handleSyncBlocks(msg *Message) {
	defer rwMutex.Unlock()
	rwMutex.Lock()

	if msg.MyHeight > uint64(len(BlockChain)) {
		var hash string
		var blocks []Block
		var is_getting_longer_chain bool
		json.Unmarshal(msg.Msg["hash_of_latest_block"], &hash)
		json.Unmarshal(msg.Msg["sync_blocks"], &blocks)
		json.Unmarshal(msg.Msg["is_getting_longer_chain"], &is_getting_longer_chain)
		BlockChain = blocks
	}

	if msg.MyHeight > uint64(len(BlockChain)) {
		var hash string
		var blocks []Block
		var isGetingLongerChain bool
		var has = false
		json.Unmarshal(msg.Msg["hash_of_latest_block"], &hash)
		json.Unmarshal(msg.Msg["sync_blocks"], &blocks)
		json.Unmarshal(msg.Msg["is_getting_longer_chain"], &isGetingLongerChain)

		count := msg.MyHeight - (uint64)(len(BlockChain))
		if isGetingLongerChain {
			count = uint64(len(blocks))
		}
		for i := count; i > 0; i = i - 1 {
			for _, v := range blocks {
				if v.Hash == hash {
					has = true
				}
			}
			if !has {
				log.Println("Warning: refuse to sync valid blocks")
				return
			}
			var block Block
			for _, v := range blocks {
				if v.Hash == hash {
					block = v
				}
			}
			if cryptoHash(&block) != hash {
				log.Println("Warning: refuse to sync valid blocks")
				return
			}
			for _, v := range BlockChain {
				if v.Hash == hash {
					hash = v.PrevHash
				}
			}
		}
		if isGetingLongerChain {
			genesisBlock := genesisBlock()
			if hash == cryptoHash(&genesisBlock) {
				BlockChain = blocks

				log.Printf("Updated Blockchain after cutting off a shorter fork by repacing it with a longer chain, result BlockChain: \n%v", BlockChain)

				go broadcastUpdated()
				disposeOrphanBlocks()
			} else {
				log.Println("Warning: refuse to replace BlockChain with a valid chain based on a different genesis block")
			}
		} else {
			if hash == BlockChain[len(BlockChain)-1].Hash {
				log.Printf("Updated BlockChain after receiving sync_blocks, result BlockChain: \n%v", BlockChain)

				go broadcastUpdated()
				disposeOrphanBlocks()
			} else {
				log.Println("Maybe my BlockChain contains a shorter fork, trying to sync")
				syncGetLongerChain(msg.From)
			}
		}
	}
}

func syncGetBlocks() {
	msg := new(Message)
	msg.From = getMyAddress()
	msg.Type = "SYNC_GET_BLOCKS"
	msg.MyHeight = uint64(len(BlockChain))
	msg.Msg = make(map[string][]byte)

	for to, _ := range listPeers {
		msg.To = to
		send(msg)
	}
	time.Sleep(1 * time.Second)
}

func syncGetLongerChain(to string) {
	msg := new(Message)
	msg.From = getMyAddress()
	msg.To = to
	msg.Type = "SYNC_GET_BLOCKS"
	msg.MyHeight = 1
	msg.Msg = make(map[string][]byte)
	is_getting_longer_chain := true
	data, err := json.Marshal(is_getting_longer_chain)
	if err != nil {
		log.Println(err)
	}
	msg.Msg["is_getting_longer_chain"] = data

	send(msg)
}

func broadcastUpdated() {
	msg := new(Message)
	msg.From = getMyAddress()
	msg.Type = "BROADCAST_UPDATED"
	msg.MyHeight = uint64(len(BlockChain))

	for to, _ := range listPeers {
		msg.To = to
		send(msg)
	}
}

func disposeOrphanBlocks() { // To retrieve orphan blocks, and check their validities.
	for orphan_block_key, orphan_block := range orphanBlocks {
		if orphan_block_key.height <= uint64(len(BlockChain)) {
			for _, v := range BlockChain {
				if v.Height == orphan_block.Height {
					go broadcastBlock(&orphan_block)
				}
			}
			delete(orphanBlocks, orphan_block_key)
		} else if orphan_block_key.height == uint64(len(BlockChain)+1) {
			if orphan_block.PrevHash == BlockChain[len(BlockChain)-1].Hash { // Valid successor
				BlockChain = append(BlockChain, orphan_block)
				go broadcastBlock(&orphan_block)
			}
			delete(orphanBlocks, orphan_block_key)
		}
	}
}
