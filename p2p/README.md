#### p2p
* go 版本 >= 1.11
* cd 到p2p目录下
* 第一个终端运行 `go run p2p.go -secio -l 10000`
* 第一个终端会打印一个链接该终端的命令,复制并粘贴到你的第二个终端,例如"go run main.go -l 10001 -d /ip4/127.0.0.1/tcp/10000/ipfs/QmUX3pYWV3R1s67bfF2q4SKtH9MxYZU6qd5QqgbyRFYzc7 -secio"
* 第二终端也有类似的命令,粘贴到第三个终端运行
* 所有终端会同步链上数据
* Value 为终端输入信息

#### 参考
[https://github.com/mycoralhealth/blockchain-tutorial/blob/master/p2p](https://github.com/mycoralhealth/blockchain-tutorial/tree/master/p2p)