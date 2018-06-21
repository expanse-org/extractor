# Loopring Extractor 中文文档

## 简介
Loopring Extractor(解析器)，遍历以太坊block及transaction将路印中继支持的合约事件及方法<br>
从交易中提取出来，并转换成中继使用的数据类型，最后使用kafka消息队列发送给中继及旷工。

## 工作流程
* 确定起始块--解析器启动后首先根据配置文件参数及数据库存储的block数据确定起始解析blockNumber
* 获取节点数据--从起始块开始遍历以太坊节点block,并批量获取transaction&transactionReceipt
* 解析事件及方法--根据合约method/event数据结构解析transaction.Input&transactionReceipt.logs,并转换成中继及旷工需要的数据结构
* 分叉检测--根据块号及parent hash判断是否有分叉,如果有分叉,生成中继/矿工支持的分叉通知数据类型
* kafka消息队列--将解析的数据及分叉数据使用kfaka消息队列发送出去

## 重要参数列表

| 参数         | 释义         |
|--------------|------------|
| log.output_paths| log输出(如果使用docker运行,需设置为/opt/loopring/extractor/logs/zap.log)|
| log.error_output_paths| err输出(如果使用docker运行,需设置为/opt/loopring/extractor/logs/err.log)|
|accessor.raw_urls|以太坊节点列表|
|accessor.fetch_tx_retry_count|获取transaction数据不成功时重试次数|
|extractor.start_block_number|解析起始块(第一次运行extractor默认值,后续使用mysql数据)|
|extractor.end_block_number|解析终止块(非debug模式下为0)|
|extractor.confirm_block_number|延迟确认块|
|extractor.debug|debug模式开关(控制部分log信息及endBlockNumber,非debug模式下为false)|
|loopring_protocol.implAbi|路印协议impl abi|
|loopring_protocol.delegateAbi|路印协议delegate abi|
|loopring_protocol.tokenRegistryAbi|路印协议token registry abi|
|loopring_protocol.address|合约地址map|
|market.token_file|中继支持的代币列表文件|
|zk_lock.zk_servers|zookeeper服务节点地址|
|kafka.brokers|kafka broker列表|

## 编译
从github上拉取代码后,运行
```bash
cd $GOPATH/src/github.com/Loopring/extractor
go build -ldflags -s -v  -o build/bin/extractor cmd/main.go
```
将在项目build/bin目录下生成extractor可执行文件

## 运行
```bash
extractor --config=config_file
```

## 依赖
* mysql数据库
* redis缓存
* 以太坊节点(集群)
* zookeeper-kafka消息队列

## 部署
* 可执行文件部署-- 环境:gcc, golang(v1.9.0以上), <br>
    部署前请根据配置文件mysql,redis,kafka,zk相关端口进行telnet测试,确保这些依赖能正常访问
* [docker](https://loopring.github.io/extractor/DOCKER_CN)

## 支持
请访问官方网站获取联系方式，获得帮助: https://loopring.org