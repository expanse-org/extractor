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

| parameter         | description         |
|--------------|------------|
| log.output_paths| N|
| log.error_output_paths| N|
|accessor.raw_urls|N|
|accessor.fetch_tx_retry_count|N|
|extractor.start_block_number|N|
|extractor.end_block_number|N|
|extractor.confirm_block_number|N|
|extractor.debug|N|
|loopring_protocol.implAbi|N|
|loopring_protocol.delegateAbi|N|
|loopring_protocol.tokenRegistryAbi|N|
|loopring_protocol.address|N|
|market.token_file|N|
|zk_lock.zk_servers|N|
|kafka_producer.brokers|N|
|kafka_consumer.brokers|N|

## Build
```bash
cd $GOPATH/src/github.com/Loopring/extractor
go build -ldflags -s -v  -o build/bin/extractor cmd/main.go
```

## Run
```bash
extractor --config=config_file
```

## Deploy
* 
* [Docker](https://loopring.github.io/extractor/DOCKER_EN)

## Support
