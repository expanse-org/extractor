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
	"github.com/Loopring/relay-lib/eventemitter"
)

func getTopic(name string) string {
	var topic string

	switch name {
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
	default:
		topic = ""
	}

	return topic
}

func Emit(name string, event interface{}) error {
	topic := getTopic(name)
	if topic == "" {
		return fmt.Errorf("emit topic is empty")
	}

	eventemitter.Emit(topic, event)
	return nil
}
