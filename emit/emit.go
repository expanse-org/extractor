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

package emit

import (
	"fmt"
	"github.com/Loopring/relay-lib/eth/contract"
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
		topic = kafka.AddressAuthorized

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

const (
	ZKNAME_EXTRACTOR = "extractor"
)

func Produce(topic string, event interface{}) error {
	zklock.TryLock(ZKNAME_EXTRACTOR)

	if topic == "" {
		return fmt.Errorf("emit topic is empty")
	}

	// todo 对接kafka
	//kafka.Produce(topic, event)
	return nil
}

func Consume() error {
	return nil
}
