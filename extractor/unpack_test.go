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

package extractor_test

import (
	"github.com/Loopring/relay-lib/eth/abi"
	"github.com/Loopring/relay-lib/eth/contract"
	"github.com/Loopring/relay-lib/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"testing"
)

const (
	erc20AbiStr         = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
	wethAbiStr          = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"guy\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]"
	implAbiStr          = "[{\"constant\":true,\"inputs\":[],\"name\":\"MARGIN_SPLIT_PERCENTAGE_BASE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ringIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"RATE_RATIO_SCALE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lrcTokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenRegistryAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegateAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderOwner\",\"type\":\"address\"},{\"name\":\"token1\",\"type\":\"address\"},{\"name\":\"token2\",\"type\":\"address\"}],\"name\":\"getTradingPairCutoffs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token1\",\"type\":\"address\"},{\"name\":\"token2\",\"type\":\"address\"},{\"name\":\"cutoff\",\"type\":\"uint256\"}],\"name\":\"cancelAllOrdersByTradingPair\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addresses\",\"type\":\"address[5]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6]\"},{\"name\":\"buyNoMoreThanAmountB\",\"type\":\"bool\"},{\"name\":\"marginSplitPercentage\",\"type\":\"uint8\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_RING_SIZE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"cutoff\",\"type\":\"uint256\"}],\"name\":\"cancelAllOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rateRatioCVSThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addressList\",\"type\":\"address[4][]\"},{\"name\":\"uintArgsList\",\"type\":\"uint256[6][]\"},{\"name\":\"uint8ArgsList\",\"type\":\"uint8[1][]\"},{\"name\":\"buyNoMoreThanAmountBList\",\"type\":\"bool[]\"},{\"name\":\"vList\",\"type\":\"uint8[]\"},{\"name\":\"rList\",\"type\":\"bytes32[]\"},{\"name\":\"sList\",\"type\":\"bytes32[]\"},{\"name\":\"feeRecipient\",\"type\":\"address\"},{\"name\":\"feeSelections\",\"type\":\"uint16\"}],\"name\":\"submitRing\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"walletSplitPercentage\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_ringIndex\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"_ringHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_miner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_feeRecipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_orderInfoList\",\"type\":\"bytes32[]\"}],\"name\":\"RingMined\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_amountCancelled\",\"type\":\"uint256\"}],\"name\":\"OrderCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_cutoff\",\"type\":\"uint256\"}],\"name\":\"AllOrdersCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_token1\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_token2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_cutoff\",\"type\":\"uint256\"}],\"name\":\"OrdersCancelled\",\"type\":\"event\"}]"
	tokenRegistryAbiStr = "[{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"unregisterToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getAddressBySymbol\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addressList\",\"type\":\"address[]\"}],\"name\":\"areAllTokensRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isTokenRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"getTokens\",\"outputs\":[{\"name\":\"addressList\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"registerToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"addresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"isTokenRegisteredBySymbol\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"TokenRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"TokenUnregistered\",\"type\":\"event\"}]"
	delegateAbiStr      = "[{\"constant\":true,\"inputs\":[{\"name\":\"owners\",\"type\":\"address[]\"},{\"name\":\"tradingPairs\",\"type\":\"bytes20[]\"},{\"name\":\"validSince\",\"type\":\"uint256[]\"}],\"name\":\"checkCutoffsBatch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"resume\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"getLatestAuthorizedAddresses\",\"outputs\":[{\"name\":\"addresses\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"cancelOrFillAmount\",\"type\":\"uint256\"}],\"name\":\"addCancelledOrFilled\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"lrcTokenAddress\",\"type\":\"address\"},{\"name\":\"miner\",\"type\":\"address\"},{\"name\":\"feeRecipient\",\"type\":\"address\"},{\"name\":\"walletSplitPercentage\",\"type\":\"uint8\"},{\"name\":\"batch\",\"type\":\"bytes32[]\"}],\"name\":\"batchTransferToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"authorizeAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenPair\",\"type\":\"bytes20\"},{\"name\":\"t\",\"type\":\"uint256\"}],\"name\":\"setTradingPairCutoffs\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelledOrFilled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"suspended\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"batch\",\"type\":\"bytes32[]\"}],\"name\":\"batchAddCancelledOrFilled\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"bytes20\"}],\"name\":\"tradingPairCutoffs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"cancelAmount\",\"type\":\"uint256\"}],\"name\":\"addCancelled\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"addressInfos\",\"outputs\":[{\"name\":\"previous\",\"type\":\"address\"},{\"name\":\"index\",\"type\":\"uint32\"},{\"name\":\"authorized\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAddressAuthorized\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"cutoffs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"suspend\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"deauthorizeAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"t\",\"type\":\"uint256\"}],\"name\":\"setCutoffs\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"number\",\"type\":\"uint32\"}],\"name\":\"AddressAuthorized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"number\",\"type\":\"uint32\"}],\"name\":\"AddressDeauthorized\",\"type\":\"event\"}]"

	protocolStr = "0xad52ab3c0e02cc3a3651ed30cfe83c3979567796"
)

var (
	erc20Abi         *abi.ABI
	wethAbi          *abi.ABI
	implAbi          *abi.ABI
	tokenRegistryAbi *abi.ABI
	delegateAbi      *abi.ABI
	protocol         common.Address
)

func init() {
	erc20Abi, _ = abi.New(erc20AbiStr)
	wethAbi, _ = abi.New(wethAbiStr)
	implAbi, _ = abi.New(implAbiStr)
	tokenRegistryAbi, _ = abi.New(tokenRegistryAbiStr)
	delegateAbi, _ = abi.New(delegateAbiStr)
	protocol = common.HexToAddress(protocolStr)
}

func TestExtractorServiceImpl_UnpackSubmitRingMethod(t *testing.T) {
	input := "0xe78aadb20000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000003e0000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000004a0000000000000000000000000000000000000000000000000000000000000054000000000000000000000000000000000000000000000000000000000000005e00000000000000000000000004bad3053d574cd54513babe21db3f09bea1d387d00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000b1018949b241d76a1ab2094f473e9befeabb5ead000000000000000000000000ae79693db742d72576db8349142f9cd8b9d85355000000000000000000000000b1018949b241d76a1ab2094f473e9befeabb5ead00000000000000000000000047fe1648b80fa04584241781488ce4c0aaca23e40000000000000000000000001b978a1d302335a6f2ebe4b8823b5e17c3c84135000000000000000000000000f079e0612e869197c5f4c7d0a95df570b163232b0000000000000000000000001b978a1d302335a6f2ebe4b8823b5e17c3c8413500000000000000000000000047fe1648b80fa04584241781488ce4c0aaca23e4000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000001043561a8829300000000000000000000000000000000000000000000000000000016345785d8a0000000000000000000000000000000000000000000000000000000000005b0277ac000000000000000000000000000000000000000000000000000000005b864dac00000000000000000000000000000000000000000000000029a2241af62c000000000000000000000000000000000000000000000000001043561a8829300000000000000000000000000000000000000000000000000000016345785d8a000000000000000000000000000000000000000000000000001043561a8829300000000000000000000000000000000000000000000000000000000000005b0277ac000000000000000000000000000000000000000000000000000000005b864dac0000000000000000000000000000000000000000000000004563918244f40000000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000001b000000000000000000000000000000000000000000000000000000000000001b000000000000000000000000000000000000000000000000000000000000001b0000000000000000000000000000000000000000000000000000000000000004e152e9cb3819f60351d9dcbdc9bd17e18e23dac72832fd0bc86cbe11a5b5a321c0bc5dffe168a513b6e776a05c0c564001f6e91f264dcf58fe48d0b9929a9de2351a24b4629480d1bdc630b52e972085cc354538310fb1d1a7299d974e6f6890351a24b4629480d1bdc630b52e972085cc354538310fb1d1a7299d974e6f689000000000000000000000000000000000000000000000000000000000000000042cb99ebec21a774ba58271b61743b913784a3e466c3425ca777fbc730368b8df1e72ab1891367ffb5bd81734df5439553992ed05e641d0af06118e0663a417c70cb200d2852e40edf741bf926764d6e5d2bdd8a82877786e24e150a6fd9632b20cb200d2852e40edf741bf926764d6e5d2bdd8a82877786e24e150a6fd9632b2"

	var ring contract.SubmitRingMethodInputs
	ring.Protocol = protocol

	data := hexutil.MustDecode("0x" + input[10:])
	value := [][]byte{}

	if err := implAbi.UnpackMethod(&ring, "submitRing", data, value); err != nil {
		t.Fatalf(err.Error())
	}

	event, err := ring.ConvertDown()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for k, v := range event.OrderList {
		t.Log(k, "protocol", v.Protocol.Hex())
		t.Log(k, "tokenS", v.TokenS.Hex())
		t.Log(k, "tokenB", v.TokenB.Hex())

		t.Log(k, "amountS", v.AmountS.String())
		t.Log(k, "amountB", v.AmountB.String())
		t.Log(k, "validSince", v.ValidSince.String())
		t.Log(k, "validUntil", v.ValidUntil.String())
		t.Log(k, "lrcFee", v.LrcFee.String())
		t.Log(k, "rateAmountS", ring.UintArgsList[k][5].String())

		t.Log(k, "marginSplitpercentage", v.MarginSplitPercentage)
		t.Log(k, "feeSelectionList", ring.Uint8ArgsList[k][0])

		t.Log(k, "buyNoMoreThanAmountB", v.BuyNoMoreThanAmountB)

		t.Log(k, "v", v.V)
		t.Log(k, "s", v.S.Hex())
		t.Log(k, "r", v.R.Hex())
	}

	t.Log("feeReceipt", event.FeeReceipt.Hex())
	t.Log("feeSelection", event.FeeSelection)
}

func TestExtractorServiceImpl_UnpackWethWithdrawalMethod(t *testing.T) {
	input := "0x2e1a7d4d0000000000000000000000000000000000000000000000000000000000000064"

	var withdrawal contract.WethWithdrawalMethod

	data := hexutil.MustDecode("0x" + input[10:])

	if err := wethAbi.UnpackMethod(&withdrawal, "withdraw", data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	}

	evt := withdrawal.ConvertDown()
	t.Logf("withdrawal event value:%s", evt.Amount)
}

func TestExtractorServiceImpl_UnpackCancelOrderMethod(t *testing.T) {
	input := "0x8c59f7ca000000000000000000000000b1018949b241d76a1ab2094f473e9befeabb5ead000000000000000000000000480037780d0b0e766941b8c5e99e685bf8812c39000000000000000000000000f079e0612e869197c5f4c7d0a95df570b163232b000000000000000000000000b1018949b241d76a1ab2094f473e9befeabb5ead00000000000000000000000047fe1648b80fa04584241781488ce4c0aaca23e400000000000000000000000000000000000000000000003635c9adc5dea00000000000000000000000000000000000000000000000000000016345785d8a0000000000000000000000000000000000000000000000000000000000005ad8a62f000000000000000000000000000000000000000000000000000000005b5c7c2f00000000000000000000000000000000000000000000000029a2241af62c00000000000000000000000000000000000000000000000000001bc16d674ec8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001b39026cca9b4e4e42ac957182e6bbeebd88d327c9368f905620b8edbf2be687af12e190eb0ec2fc5b337487834aeb9ce9df2f0275f281b3e7ca5bdec13246444f"

	var method contract.CancelOrderMethod

	data := hexutil.MustDecode("0x" + input[10:])

	//for i := 0; i < len(data)/32; i++ {
	//	t.Logf("index:%d -> %s", i, common.ToHex(data[i*32:(i+1)*32]))
	//}

	if err := implAbi.UnpackMethod(&method, "cancelOrder", data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	}

	order, cancelAmount, err := method.ConvertDown()
	if err != nil {
		t.Fatalf(err.Error())
	}

	order.DelegateAddress = common.HexToAddress("0xf49733091a3e1ddec740bca4c325f8aaee6ee307")
	order.Hash = order.GenerateHash()
	t.Log("delegate", order.DelegateAddress.Hex())
	t.Log("orderHash", order.Hash.Hex())
	t.Log("owner", order.Owner.Hex())
	t.Log("wallet", order.WalletAddress.Hex())
	t.Log("auth", order.AuthAddr.Hex())
	t.Log("tokenS", order.TokenS.Hex())
	t.Log("tokenB", order.TokenB.Hex())
	t.Log("amountS", order.AmountS.String())
	t.Log("amountB", order.AmountB.String())
	t.Log("validSince", order.ValidSince.String())
	t.Log("validUntil", order.ValidUntil.String())
	t.Log("lrcFee", order.LrcFee.String())
	t.Log("cancelAmount", method.OrderValues[5].String())
	t.Log("buyNoMoreThanAmountB", order.BuyNoMoreThanAmountB)
	t.Log("marginSplitpercentage", order.MarginSplitPercentage)
	t.Log("v", order.V)
	t.Log("s", order.S.Hex())
	t.Log("r", order.R.Hex())
	t.Log("cancelAmount", cancelAmount)
}

func TestExtractorServiceImpl_UnpackApproveMethod(t *testing.T) {
	input := "0x095ea7b300000000000000000000000045aa504eb94077eec4bf95a10095a8e3196fc5910000000000000000000000000000000000000000000000008ac7230489e80000"

	var method contract.ApproveMethod

	data := hexutil.MustDecode("0x" + input[10:])
	for i := 0; i < len(data)/32; i++ {
		t.Logf("index:%d -> %s", i, common.ToHex(data[i*32:(i+1)*32]))
	}

	if err := erc20Abi.UnpackMethod(&method, "approve", data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	}

	approve := method.ConvertDown()
	t.Logf("approve spender:%s, value:%s", approve.Spender.Hex(), approve.Amount.String())
}

func TestExtractorServiceImpl_UnpackTransferMethod(t *testing.T) {
	input := "0xa9059cbb0000000000000000000000008311804426a24495bd4306daf5f595a443a52e32000000000000000000000000000000000000000000000000000000174876e800"
	data := hexutil.MustDecode("0x" + input[10:])

	var method contract.TransferMethod
	if err := erc20Abi.UnpackMethod(&method, "transfer", data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	}
	transfer := method.ConvertDown()

	t.Logf("transfer receiver:%s, value:%s", transfer.Receiver.Hex(), transfer.Amount.String())
}

func TestExtractorServiceImpl_UnpackTransferEvent(t *testing.T) {
	inputs := []string{
		"0x00000000000000000000000000000000000000000000001d2666491321fc5651",
		"0x0000000000000000000000000000000000000000000000008ac7230489e80000",
		"0x0000000000000000000000000000000000000000000000004c0303a413a39039",
		"0x000000000000000000000000000000000000000000000000016345785d8a0000",
	}
	transfer := &contract.TransferEvent{}

	for _, input := range inputs {
		data := hexutil.MustDecode(input)

		if err := erc20Abi.UnpackEvent(transfer, "Transfer", []byte{}, [][]byte{data}); err != nil {
			t.Fatalf(err.Error())
		}

		t.Logf("transfer value:%s", transfer.Value.String())
	}
}

func TestExtractorServiceImpl_UnpackRingMinedEvent(t *testing.T) {
	input := "0x00000000000000000000000000000000000000000000000000000000000000070000000000000000000000004bad3053d574cd54513babe21db3f09bea1d387d0000000000000000000000004bad3053d574cd54513babe21db3f09bea1d387d0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000eece69a21bb35f7566d4d7e447cb2765cf464c308ba0352d6ad90af4a744794eb0000000000000000000000001b978a1d302335a6f2ebe4b8823b5e17c3c84135000000000000000000000000f079e0612e869197c5f4c7d0a95df570b163232b000000000000000000000000000000000000000000000000016345785d8a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffba9c6e7dbb0c00006987b1498573ad4fed2d2a1becb054c57d351f775c1dd3d80a42a25dd31c18e3000000000000000000000000b1018949b241d76a1ab2094f473e9befeabb5ead000000000000000000000000ae79693db742d72576db8349142f9cd8b9d8535500000000000000000000000000000000000000000000001db12d6c17abe45651000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016cdb44ad2b111aa0000000000000000000000000000000000000000000000000000000000000000"
	//input := "0x00000000000000000000000000000000000000000000000000000000000000080000000000000000000000004bad3053d574cd54513babe21db3f09bea1d387d0000000000000000000000004bad3053d574cd54513babe21db3f09bea1d387d0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000e779a662897f805cee228e4c0349ec8a5c05c190652287b47daddc3008d78a28b000000000000000000000000b1018949b241d76a1ab2094f473e9befeabb5ead000000000000000000000000ae79693db742d72576db8349142f9cd8b9d8535500000000000000000000000000000000000000000000001043561a8829300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016cdb44ad2b111b40000000000000000000000000000000000000000000000000000000000000000af78d9d04c29924ff9dcdda4f034f77e230d186415fe433bc653e980d4d6771f0000000000000000000000001b978a1d302335a6f2ebe4b8823b5e17c3c84135000000000000000000000000f079e0612e869197c5f4c7d0a95df570b163232b00000000000000000000000000000000000000000000000000c297138f8e6f8100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffba9c6e7dbb0c0000"
	ringmined := &contract.RingMinedEvent{}

	data := hexutil.MustDecode(input)

	for i := 0; i < len(data)/32; i++ {
		t.Logf("index:%d -> %s", i, common.ToHex(data[i*32:(i+1)*32]))
	}

	if err := implAbi.UnpackEvent(ringmined, "RingMined", data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	}

	evt, fills, err := ringmined.ConvertDown()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for k, fill := range fills {
		t.Logf("k:%d --> ringindex:%s", k, fill.RingIndex.String())
		t.Logf("k:%d --> fillIndex:%s", k, fill.FillIndex.String())
		t.Logf("k:%d --> orderhash:%s", k, fill.OrderHash.Hex())
		t.Logf("k:%d --> preorder:%s", k, fill.PreOrderHash.Hex())
		t.Logf("k:%d --> nextorder:%s", k, fill.NextOrderHash.Hex())
		t.Logf("k:%d --> owner:%s", k, fill.Owner.Hex())
		t.Logf("k:%d --> tokenS:%s", k, fill.TokenS.Hex())
		t.Logf("k:%d --> tokenB:%s", k, fill.TokenB.Hex())
		t.Logf("k:%d --> amountS:%s", k, fill.AmountS.String())
		t.Logf("k:%d --> amountB:%s", k, fill.AmountB.String())
		t.Logf("k:%d --> lrcReward:%s", k, fill.LrcReward.String())
		t.Logf("k:%d --> lrcFee:%s", k, fill.LrcFee.String())
		t.Logf("k:%d --> splitS:%s", k, fill.SplitS.String())
		t.Logf("k:%d --> splitB:%s", k, fill.SplitB.String())
	}

	t.Logf("totalLrcFee:%s", evt.TotalLrcFee.String())
	t.Logf("tradeAmount:%d", evt.TradeAmount)
}

func TestExtractorServiceImpl_UnpackDepositEvent(t *testing.T) {
	input := "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000"
	deposit := &contract.WethDepositEvent{}

	data := hexutil.MustDecode(input)

	if err := wethAbi.UnpackEvent(deposit, "Deposit", data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	} else {
		t.Logf("deposit value:%s", deposit.Value.String())
	}
}

func TestExtractorServiceImpl_UnpackTokenRegistryEvent(t *testing.T) {
	input := "0x000000000000000000000000f079e0612e869197c5f4c7d0a95df570b163232b0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000457455448"

	tokenRegistry := &contract.TokenRegisteredEvent{}

	data := hexutil.MustDecode(input)

	if err := tokenRegistryAbi.UnpackEvent(tokenRegistry, contract.EVENT_TOKEN_REGISTERED, data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	} else {
		t.Logf("TokenRegistered symbol:%s, address:%s", tokenRegistry.Symbol, tokenRegistry.Token.Hex())
	}
}

func TestExtractorServiceImpl_UnpackTokenUnRegistryEvent(t *testing.T) {
	input := "0x000000000000000000000000529540ee6862158f47d647ae023098f6705210a90000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000457455448"

	tokenUnRegistry := &contract.TokenUnRegisteredEvent{}

	data := hexutil.MustDecode(input)

	if err := tokenRegistryAbi.UnpackEvent(tokenUnRegistry, "TokenUnregistered", data, [][]byte{}); err != nil {
		t.Fatalf(err.Error())
	} else {
		t.Logf("TokenUnregistered symbol:%s, address:%s", tokenUnRegistry.Symbol, tokenUnRegistry.Token.Hex())
	}
}

func TestExtractorServiceImpl_Compare(t *testing.T) {
	str1 := "547722557505166136913"
	str2 := "1000000000000000000000"
	num1, _ := big.NewInt(0).SetString(str1, 0)
	num2, _ := big.NewInt(0).SetString(str2, 0)
	if num1.Cmp(num2) > 0 {
		t.Logf("%s > %s", str1, str2)
	} else {
		t.Logf("%s <= %s", str1, str2)
	}
}

func TestExtractorServiceImpl_UnpackNumbers(t *testing.T) {
	str1 := "0xffffffffffffffffffffffffffffffffffffffffffffffffffa1d2c1fb1c2d9f"
	str2 := "0xffffffffffffffffffffffffffffffffffffffffffffffffff90c5f64e557fa4"
	str3 := "0x0000000000000000000000000000000000000000000000026508392204063330"
	str4 := "0x0000000000000000000000000000000000000000000000031307535724740700"
	list := []string{str1, str2, str3, str4}

	for _, v := range list {
		n1 := safeBig(v)
		t.Logf("init data:%s -> number:%s", v, n1.String())
	}
}

func safeBig(input string) *big.Int {
	bytes := hexutil.MustDecode(input)
	num := new(big.Int).SetBytes(bytes[:])
	if bytes[0] > uint8(128) {
		num.Xor(types.MaxUint256, num)
		num.Not(num)
	}
	return num
}
