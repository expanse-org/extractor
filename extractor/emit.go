/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package extractor

import (
	"fmt"
	"github.com/Loopring/relay-lib/eth/contract"
	"github.com/Loopring/relay-lib/eth/types"
	"github.com/Loopring/relay-lib/kafka"
	"github.com/Loopring/relay-lib/zklock"
)

func Topic(name string) string {
	var topic string

	switch name {
	// methods
	case contract.METHOD_SUBMIT_RING:
		topic = kafka.Miner_SubmitRing_Method

	case contract.METHOD_CANCEL_ORDER:
		topic = kafka.CancelOrder

	case contract.METHOD_CUTOFF_ALL:
		topic = kafka.CutoffAll

	case contract.METHOD_CUTOFF_PAIR:
		topic = kafka.CutoffPair

	case contract.METHOD_APPROVE:
		topic = kafka.Approve

	case contract.METHOD_TRANSFER:
		topic = kafka.Transfer

	case contract.METHOD_WETH_DEPOSIT:
		topic = kafka.WethDeposit

	case contract.METHOD_WETH_WITHDRAWAL:
		topic = kafka.WethWithdrawal

	// events
	case contract.EVENT_ORDER_CANCELLED:
		topic = kafka.CancelOrder

	case contract.EVENT_CUTOFF_ALL:
		topic = kafka.CutoffAll

	case contract.EVENT_CUTOFF_PAIR:
		topic = kafka.CutoffPair

	case contract.EVENT_TRANSFER:
		topic = kafka.Transfer

	case contract.EVENT_APPROVAL:
		topic = kafka.Approve

	case contract.EVENT_WETH_DEPOSIT:
		topic = kafka.WethDeposit

	case contract.EVENT_WETH_WITHDRAWAL:
		topic = kafka.WethWithdrawal

	case contract.EVENT_TOKEN_REGISTERED:
		topic = kafka.TokenRegistered

	case contract.EVENT_TOKEN_UNREGISTERED:
		topic = kafka.TokenUnRegistered

	case contract.EVENT_ADDRESS_AUTHORIZED:
		topic = kafka.AddressAuthorized

	case contract.EVENT_ADDRESS_DEAUTHORIZED:
		topic = kafka.AddressDeAuthorized

	default:
		topic = ""
	}

	return topic
}

func RingMinedTopic(isFill bool) string {
	if isFill {
		return kafka.OrderFilled
	}
	return kafka.RingMined
}

func EthTxTopic(isTransfer bool) string {
	if isTransfer {
		return kafka.EthTransfer
	}
	return kafka.UnsupportedContract
}

func NewBlockTopic(isFinished bool) string {
	if isFinished {
		return kafka.Block_End
	}
	return kafka.Block_New
}

func ForkEventTopic() string { return kafka.ExtractorFork }

func SyncCompleteTopic() string { return kafka.SyncChainComplete }

const (
	ZKNAME_EXTRACTOR     = "extractor"
	KAFKA_CONSUMER_TOPIC = kafka.PendingTransaction
	KAFKA_CONSUMER_GROUP = "extractor_pending_transaction"
)

var (
	producer *kafka.MessageProducer
	register *kafka.ConsumerRegister
)

func RegistryEmitter(zkOpt zklock.ZkLockConfig, producerOpt, consumerOpt kafka.KafkaOptions, service ExtractorService) error {
	if _, err := zklock.Initialize(zkOpt); err != nil {
		return err
	}

	brokers := producerOpt.Brokers
	producer = &kafka.MessageProducer{}
	if err := producer.Initialize(brokers); err != nil {
		return err
	}

	if len(consumerOpt.Brokers) < 1 {
		return fmt.Errorf("kafka consumer brokers should not be empty")
	}
	register = &kafka.ConsumerRegister{}
	register.Initialize(consumerOpt.Brokers[0])
	if err := register.RegisterTopicAndHandler(KAFKA_CONSUMER_TOPIC, KAFKA_CONSUMER_GROUP, types.Transaction{}, service.WatchingPendingTransaction); err != nil {
		return err
	}

	return nil
}

// todo:这里貌似释放非常慢
func UnRegistryEmitter() {
	zklock.ReleaseLock(ZKNAME_EXTRACTOR)
	producer.Close()
	register.Close()
}

func Produce(topic string, event interface{}) error {
	zklock.TryLock(ZKNAME_EXTRACTOR)

	if topic == "" {
		return fmt.Errorf("emit topic is empty")
	}

	// todo 对接kafka
	producer.SendMessage(topic, event, nextKey(topic))
	return nil
}

// todo
func nextKey(topic string) string {
	return ""
}
