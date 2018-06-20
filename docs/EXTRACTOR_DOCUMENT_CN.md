# Loopring Extractor 中文文档

## 简介
Loopring Extractor(解析器)是路印技术生态重要组成部分，遍历以太坊block及transaction，<br>
将路印中继支持的合约事件及方法从交易中提取出来，并转换成中继使用的数据类型，最后<br>
使用kafka消息队列发送给中继及旷工。

## 依赖
* mysql数据库
* redis缓存
* 以太坊节点(集群)
* zookeeper-kafka消息队列

## 工作流程

* 确定起始块--解析器启动后首先根据配置文件参数及数据库存储的block数据确定起始解析blockNumber
* 获取节点数据--从起始块开始遍历以太坊节点block,并批量获取transaction&transactionReceipt
* 解析事件及方法--根据合约method/event数据结构解析transaction.Input&transactionReceipt.logs,并转换成中继及旷工需要的数据结构
* 分叉检测--根据块号及parent hash判断是否有分叉,如果有分叉,生成中继/矿工支持的分叉通知数据类型
* kafka消息队列--将解析的数据及分叉数据使用kfaka消息队列发送出去

---

## 编译
从github上拉取代码后,运行
```bash
cd $GOPATH/src/github.com/Loopring/extractor
go build -ldflags -s -v  -o build/bin/extractor cmd/main.go
```
将在项目build/bin目录下生成extractor可执行文件

## 运行
运行
```bash
extractor --config=config_file
```

## 支持
请访问官方网站获取联系方式，获得帮助: https://loopring.org