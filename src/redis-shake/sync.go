// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package run

import (
	"golang.org/x/sync/semaphore"

	"github.com/alibaba/RedisShake/pkg/libs/log"

	"github.com/alibaba/RedisShake/redis-shake/common"
	"github.com/alibaba/RedisShake/redis-shake/configure"
	"github.com/alibaba/RedisShake/redis-shake/dbSync"
)

// main struct
type CmdSync struct {
	dbSyncers []*dbSync.DbSyncer
}

// return send buffer length, delay channel length, target db offset
func (cmd *CmdSync) GetDetailedInfo() interface{} {
	ret := make([]map[string]interface{}, len(cmd.dbSyncers))
	for i, syncer := range cmd.dbSyncers {
		if syncer == nil {
			continue
		}
		ret[i] = syncer.GetExtraInfo()
	}
	return ret
}

func (cmd *CmdSync) Main() {
	type syncNode struct {
		id                int
		source            string
		sourcePassword    string
		target            []string
		targetPassword    string
		slotLeftBoundary  int
		slotRightBoundary int
	}

	var slotDistribution []utils.SlotOwner
	var err error
	if conf.Options.SourceType == conf.RedisTypeCluster && conf.Options.ResumeFromBreakPoint {
		if slotDistribution, err = utils.GetSlotDistribution(conf.Options.SourceAddressList[0], conf.Options.SourceAuthType,
			conf.Options.SourcePasswordRaw, false); err != nil {
			log.Errorf("get source slot distribution failed: %v", err)
			return
		}
	}

	// source redis number
	total := utils.GetTotalLink()
	syncChan := make(chan syncNode, total)
	cmd.dbSyncers = make([]*dbSync.DbSyncer, total)
	for i, source := range conf.Options.SourceAddressList {
		var target []string
		if conf.Options.TargetType == conf.RedisTypeCluster {
			target = conf.Options.TargetAddressList
		} else {
			// round-robin pick
			pick := utils.PickTargetRoundRobin(len(conf.Options.TargetAddressList))
			target = []string{conf.Options.TargetAddressList[pick]}
		}

		// fetch slot boundary
		leftSlotBoundary, rightSlotBoundary := utils.GetSlotBoundary(slotDistribution, source)

		nd := syncNode{
			id:                i,
			source:            source,
			sourcePassword:    conf.Options.SourcePasswordRaw,
			target:            target,
			targetPassword:    conf.Options.TargetPasswordRaw,
			slotLeftBoundary:  leftSlotBoundary,
			slotRightBoundary: rightSlotBoundary,
		}
		syncChan <- nd
	}

	//var wg sync.WaitGroup
	//wg.Add(len(conf.Options.SourceAddressList))
	maxFullsyncs := semaphore.NewWeighted(int64(conf.Options.SourceRdbParallel))
	for {
		nd, ok := <-syncChan
		if !ok {
			break
		}

		// one sync link corresponding to one DbSyncer
		ds := dbSync.NewDbSyncer(nd.id, nd.source, nd.sourcePassword, nd.target, nd.targetPassword,
			nd.slotLeftBoundary, nd.slotRightBoundary, conf.Options.HttpProfile+nd.id, maxFullsyncs)
		cmd.dbSyncers[nd.id] = ds
		// run in routine
		go ds.Sync()
	}
	close(syncChan)

	// never quit because increment syncing is always running
	select {}
}
