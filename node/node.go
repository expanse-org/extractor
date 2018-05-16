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

package node

import (
	"sync"

	"github.com/Loopring/extractor/dao"
	"github.com/Loopring/extractor/extractor"
	"github.com/Loopring/relay-lib/eth/accessor"
	"github.com/Loopring/relay-lib/log"
	"go.uber.org/zap"
)

type Node struct {
	globalConfig *GlobalConfig
	rdsService   dao.RdsService
	extractor    extractor.ExtractorService

	stop   chan struct{}
	lock   sync.RWMutex
	logger *zap.Logger
}

func NewNode(logger *zap.Logger, globalConfig *GlobalConfig) *Node {
	n := &Node{}
	n.logger = logger
	n.globalConfig = globalConfig

	// register
	n.registerMysql()
	n.registerAccessor()
	n.registerExtractor()
	//cache.NewCache(n.globalConfig.Redis)

	return n
}

func (n *Node) Start() {
	n.extractor.Start()
}

func (n *Node) Wait() {
	n.lock.RLock()

	stop := n.stop
	n.lock.RUnlock()

	<-stop
}

// todo
func (n *Node) Stop() {
}

func (n *Node) registerMysql() {
	n.rdsService = dao.NewDb(&n.globalConfig.Mysql)
}

func (n *Node) registerAccessor() {
	if err := accessor.Initialize(n.globalConfig.Accessor); nil != err {
		log.Fatalf("err:%s", err.Error())
	}
}

func (n *Node) registerExtractor() {
	extractor.NewExtractorService(&n.globalConfig.Extractor, n.rdsService)
}
