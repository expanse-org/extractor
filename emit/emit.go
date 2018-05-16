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
	eventemitter "github.com/Loopring/relay-lib/kafka"
)

func Topic(name string) string {
	var topic string

	switch name {
	// methods
	case contract.METHOD_SUBMIT_RING:
		topic = eventemitter.Miner_SubmitRing_Method

	case contract.METHOD_CANCEL_ORDER:
		topic = eventemitter.CancelOrder

	case contract.METHOD_CUTOFF_ALL:
		topic = eventemitter.CutoffAll

	case contract.METHOD_CUTOFF_PAIR:
		topic = eventemitter.CutoffPair

	case contract.METHOD_APPROVE:
		topic = eventemitter.Approve

	case contract.METHOD_TRANSFER:
		topic = eventemitter.Transfer

	case contract.METHOD_WETH_DEPOSIT:
		topic = eventemitter.WethDeposit

	case contract.METHOD_WETH_WITHDRAWAL:
		topic = eventemitter.WethWithdrawal

	// events
	case contract.EVENT_ORDER_CANCELLED:
		topic = eventemitter.CancelOrder

	case contract.EVENT_CUTOFF_ALL:
		topic = eventemitter.CutoffAll

	case contract.EVENT_CUTOFF_PAIR:
		topic = eventemitter.CutoffPair

	case contract.EVENT_TRANSFER:
		topic = eventemitter.Transfer

	case contract.EVENT_APPROVAL:
		topic = eventemitter.Approve

	case contract.EVENT_WETH_DEPOSIT:
		topic = eventemitter.WethDeposit

	case contract.EVENT_WETH_WITHDRAWAL:
		topic = eventemitter.WethWithdrawal

	case contract.EVENT_TOKEN_REGISTERED:
		topic = eventemitter.TokenRegistered

	case contract.EVENT_TOKEN_UNREGISTERED:
		topic = eventemitter.TokenUnRegistered

	case contract.EVENT_ADDRESS_AUTHORIZED:
		topic = eventemitter.AddressAuthorized

	case contract.EVENT_ADDRESS_DEAUTHORIZED:
		topic = eventemitter.AddressAuthorized

	default:
		topic = ""
	}

	return topic
}

func RingMinedTopic(isFill bool) string {
	if isFill {
		return eventemitter.OrderFilled
	}
	return eventemitter.RingMined
}

func EthTxTopic(isTransfer bool) string {
	if isTransfer {
		return eventemitter.EthTransfer
	}
	return eventemitter.UnsupportedContract
}

func Emit(topic string, event interface{}) error {
	if topic == "" {
		return fmt.Errorf("emit topic is empty")
	}

	// todo 对接kafka
	//eventemitter.Emit(topic, event)
	return nil
}
