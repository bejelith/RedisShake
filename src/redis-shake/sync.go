// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package run

import (
	"github.com/alibaba/RedisShake/redis-shake/dbSync/latencymonitor"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/slot"
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
	var slotDistribution []utils.SlotOwner
	var err error
	if conf.Options.SourceType == conf.RedisTypeCluster {
		if slotDistribution, err = utils.GetSlotDistribution(conf.Options.SourceAddressList[0], conf.Options.SourceAuthType,
			conf.Options.SourcePasswordRaw, false); err != nil {
			log.Errorf("get source slot distribution failed: %v", err)
			return
		}
		latencymonitor.NewSyntheticProducer(slotDistribution, conf.Options.SourcePasswordRaw, conf.Options.SourceTLSEnable).Run()
	} else {
		for _, sourceAddress := range conf.Options.SourceAddressList {
			slotDistribution = append(slotDistribution, utils.SlotOwner{
				Master: sourceAddress,
				Slave:  []string{},
			})
		}
	}
	log.Infof("Source Slots detected: %v", slotDistribution)
	// source redis number
	total := utils.GetTotalLink()
	syncChan := make(chan slot.SyncNode, total)
	cmd.dbSyncers = make([]*dbSync.DbSyncer, total)
	for i, node := range slotDistribution {
		var target []string
		if conf.Options.TargetType == conf.RedisTypeCluster {
			target = conf.Options.TargetAddressList
		} else {
			// round-robin pick
			pick := utils.PickTargetRoundRobin(len(conf.Options.TargetAddressList))
			target = []string{conf.Options.TargetAddressList[pick]}
		}

		nd := slot.SyncNode{
			Id:                i,
			Source:            node.Master,
			SourcePassword:    conf.Options.SourcePasswordRaw,
			Slaves:            node.Slave,
			Target:            target,
			TargetPassword:    conf.Options.TargetPasswordRaw,
			SlotLeftBoundary:  node.SlotLeftBoundary,
			SlotRightBoundary: node.SlotRightBoundary,
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
		ds := dbSync.NewDbSyncer(&nd, conf.Options.HttpProfile+nd.Id, maxFullsyncs)
		cmd.dbSyncers[nd.Id] = ds
		// run in routine
		go ds.Sync()
	}
	close(syncChan)

	// never quit because increment syncing is always running
	select {}
}
