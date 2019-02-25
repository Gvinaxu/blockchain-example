## 最基本的POW实现

## 使用
main方法可以启动本地区块链网络

写入区块
-
curl -X POST  http://127.0.0.1:1234/write -H 'content-type: application/json' -d '
    {
        "value":"a"
    }
'

获取链上区块
-
curl -X GET http://127.0.0.1:1234/get 

##
具体请求如何操作 请参照`core/net_service.go` 源码很简单 就不详解了