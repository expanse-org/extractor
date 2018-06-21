# extractor

* [中文文档(Chinese version)](https://loopring.github.io/extractor/EXTRACTOR_DOCUMENT_CN)

## Introduction

loopring extractor access block and transactions from ethereum node, unpack the transaction.input<br>
and transactionReceipt.logs to data struct which used by miner&relay-cluster.<br>
at the same time, the fork detector will check and create eth node fork event.
all of the events will send to relay and miner access of kafka.


## Dependencies

* mysql
* redis
* zookeeper-kafka
* ethereum-node

## Procedure

## Config

## Build
```bash
cd $GOPATH/src/github.com/Loopring/extractor
go build -ldflags -s -v  -o build/bin/extractor cmd/main.go
```

## Run
```bash
extractor --config=config_file
```

## Docker
[English version](https://loopring.github.io/extractor/DOCKER_EN)

## Support
