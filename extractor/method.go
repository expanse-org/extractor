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
	"github.com/Loopring/relay-lib/eth/abi"
	"github.com/Loopring/relay-lib/eth/contract"
	ethtyp "github.com/Loopring/relay-lib/eth/types"
	"github.com/Loopring/relay-lib/eventemitter"
	"github.com/Loopring/relay-lib/log"
	"github.com/Loopring/relay-lib/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type MethodData struct {
	types.TxInfo
	Method interface{}
	Abi    *abi.ABI
	Id     string
	Name   string
	Input  string
}

func newMethodData(method *abi.Method, cabi *abi.ABI) MethodData {
	var c MethodData

	c.Id = common.ToHex(method.Id())
	c.Name = method.Name
	c.Abi = cabi

	return c
}

func (method *MethodData) FullFilled(tx *ethtyp.Transaction, gasUsed, blockTime *big.Int, status types.TxStatus, methodName string) {
	method.TxInfo = setTxInfo(tx, gasUsed, blockTime, methodName)
	method.Input = tx.Input
	method.TxLogIndex = 0
	method.Status = status
}

// UnpackMethod v should be ptr
func (m MethodData) handleMethod(tx *ethtyp.Transaction) error {
	var (
		event interface{}
		err   error
	)

	switch m.Name {
	case contract.METHOD_SUBMIT_RING:
		event, err = m.fullFillSubmitRing()
	}

	if err != nil {
		return err
	}

	return EmitEvent(m.EmitTopic(), event)
}

func (m MethodData) beforeUnpack() (err error) {
	switch m.Name {
	case contract.METHOD_CANCEL_ORDER:
		if m.DelegateAddress == types.NilAddress {
			err = fmt.Errorf("cancelOrder method cann't get delegate address")
		}
	}

	return err
}

func (m MethodData) afterUnpack() error {

}

func (m MethodData) unpack(tx *ethtyp.Transaction) (err error) {
	data := hexutil.MustDecode("0x" + tx.Input[10:])
	err = m.Abi.Unpack(m.Method, m.Name, data, [][]byte{})
	return err
}

func (m MethodData) fullFillSubmitRing() (event *types.SubmitRingMethodEvent, err error) {
	src, ok := m.Method.(*contract.SubmitRingMethodInputs)
	if !ok {
		return nil, fmt.Errorf("submitRing method inputs type error:%s", err.Error())
	}

	if event, err = src.ConvertDown(); err != nil {
		return event, fmt.Errorf("submitRing method inputs convert error:%s", err.Error())
	}

	// set txinfo for event
	event.TxInfo = m.TxInfo
	if event.Status == types.TX_STATUS_FAILED {
		event.Err = fmt.Errorf("method %s transaction failed", contract.METHOD_SUBMIT_RING)
	}

	// 不需要发送订单到gateway
	//for _, v := range event.OrderList {
	//	v.Hash = v.GenerateHash()
	//	log.Debugf("extractor,tx:%s submitRing method orderHash:%s,owner:%s,tokenS:%s,tokenB:%s,amountS:%s,amountB:%s", event.TxHash.Hex(), v.Hash.Hex(), v.Owner.Hex(), v.TokenS.Hex(), v.TokenB.Hex(), v.AmountS.String(), v.AmountB.String())
	//	eventemitter.Emit(eventemitter.GatewayNewOrder, v)
	//}

	log.Debugf("extractor,tx:%s submitRing method gas:%s, gasprice:%s, status:%s", event.TxHash.Hex(), event.GasUsed.String(), event.GasPrice.String(), types.StatusStr(event.Status))

	return event, nil
}

func (m MethodData) fullFillCancelOrder() (event *types.OrderCancelledEvent, err error) {
	src, ok := m.Method.(*contract.CancelOrderMethod)
	if !ok {
		return nil, fmt.Errorf("cancelOrder method inputs type error")
	}

	order, cancelAmount, _ := src.ConvertDown()
	order.Protocol = m.Protocol
	order.DelegateAddress = m.DelegateAddress
	order.Hash = order.GenerateHash()

	// 发送到txmanager
	tmCancelEvent := &types.OrderCancelledEvent{}
	tmCancelEvent.TxInfo = m.TxInfo
	tmCancelEvent.OrderHash = order.Hash
	tmCancelEvent.AmountCancelled = cancelAmount

	log.Debugf("extractor,tx:%s cancelOrder method order tokenS:%s,tokenB:%s,amountS:%s,amountB:%s", event.TxHash.Hex(), order.TokenS.Hex(), order.TokenB.Hex(), order.AmountS.String(), order.AmountB.String())

	return tmCancelEvent, nil
}

func (m MethodData) fullFillCutoffAll() (event *types.CutoffEvent, err error) {
	src, ok := m.Method.(*contract.CutoffMethod)
	if !ok {
		return nil, fmt.Errorf("cutoffAll method inputs type error")
	}

	event = src.ConvertDown()
	event.TxInfo = contract.TxInfo
	event.Owner = event.From
	log.Debugf("extractor,tx:%s cutoff method owner:%s, cutoff:%d, status:%d", event.TxHash.Hex(), event.Owner.Hex(), event.Cutoff.Int64(), event.Status)

	return event, err
}

func (m MethodData) fullFillCutoffPair() (event *types.CutoffPairEvent, err error) {
	src, ok := m.Method.(*contract.CutoffPairMethod)
	if !ok {
		return nil, fmt.Errorf("cutoffPair method inputs type error")
	}

	event = src.ConvertDown()
	event.TxInfo = m.TxInfo
	event.Owner = cutoffpair.From

	log.Debugf("extractor,tx:%s cutoffpair method owenr:%s, token1:%s, token2:%s, cutoff:%d", event.TxHash.Hex(), event.Owner.Hex(), event.Token1.Hex(), event.Token2.Hex(), event.Cutoff.Int64())

	return
}

func (m MethodData) fullFillApprove() (event *types.ApprovalEvent, err error) {
	src, ok := m.Method.(*contract.ApproveMethod)
	if !ok {
		return nil, fmt.Errorf("approve method inputs type error")
	}

	event = src.ConvertDown()
	event.TxInfo = m.TxInfo
	event.Owner = m.From

	log.Debugf("extractor,tx:%s approve method owner:%s, spender:%s, value:%s", event.TxHash.Hex(), event.Owner.Hex(), event.Spender.Hex(), event.Amount.String())

	return
}

func (m MethodData) fullFillTransfer() (event *types.TransferEvent, err error) {
	src := m.Method.(*contract.TransferMethod)

	event = src.ConvertDown()
	event.Sender = m.From
	event.TxInfo = m.TxInfo

	log.Debugf("extractor,tx:%s transfer method sender:%s, receiver:%s, value:%s", event.TxHash.Hex(), event.Sender.Hex(), event.Receiver.Hex(), event.Amount.String())

	return
}

func (m MethodData) fullFillWethDeposit() (event *types.WethDepositEvent, err error) {
	src := m.Method.(*contract.WethWithdrawalMethod)

	event.Dst = m.From
	event.Amount = m.Value
	event.TxInfo = m.TxInfo

	log.Debugf("extractor,tx:%s wethDeposit method from:%s, to:%s, value:%s", event.TxHash.Hex(), deposit.From.Hex(), deposit.To.Hex(), deposit.Amount.String())

	return
}

func (m MethodData) fullFillWethWithdrawal() (event *types.WethWithdrawalEvent, err error) {
	src := m.Method.(*contract.WethWithdrawalMethod)

	event = src.ConvertDown()
	event.Src = m.From
	event.TxInfo = m.TxInfo

	log.Debugf("extractor,tx:%s wethWithdrawal method from:%s, to:%s, value:%s", contractData.TxHash.Hex(), withdrawal.From.Hex(), withdrawal.To.Hex(), withdrawal.Amount.String())

	return
}

func (m MethodData) EmitTopic() string {
	var topic string

	switch m.Name {
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
	}

	return topic
}

func EmitEvent(topic string, event interface{}) error {
	if topic == "" {
		return fmt.Errorf("emit topic is empty")
	}

	eventemitter.Emit(topic, event)
	return nil
}
