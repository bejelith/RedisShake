package dbSync

import (
	"bufio"
	"context"
	"github.com/alibaba/RedisShake/pkg/libs/log"
	"github.com/alibaba/RedisShake/redis-shake/base"
	"github.com/alibaba/RedisShake/redis-shake/common"
	"github.com/alibaba/RedisShake/redis-shake/metric"
	"golang.org/x/sync/semaphore"
	"io"
	"strings"
	"time"

	"github.com/alibaba/RedisShake/redis-shake/checkpoint"
	"github.com/alibaba/RedisShake/redis-shake/configure"
)

// one sync link corresponding to one DbSyncer
func NewDbSyncer(id int, source, sourcePassword string, target []string, targetPassword string,
	slotLeftBoundary, slotRightBoundary int, httpPort int, maxFullsyncs *semaphore.Weighted) *DbSyncer {
	ds := &DbSyncer{
		id:                         id,
		source:                     source,
		sourcePassword:             sourcePassword,
		target:                     target,
		targetPassword:             targetPassword,
		slotLeftBoundary:           slotLeftBoundary,
		slotRightBoundary:          slotRightBoundary,
		httpProfilePort:            httpPort,
		enableResumeFromBreakPoint: conf.Options.ResumeFromBreakPoint,
		checkpointName:             utils.CheckpointKey, // default, may be modified
		WaitFull:                   make(chan struct{}),
		Restart:                    make(chan struct{}),
		MaxParallelFullSyncs:       maxFullsyncs,
		fullSyncRetryCounter:       0,
	}

	// add metric
	metric.AddMetric(id)

	return ds
}

type DbSyncer struct {
	id int // current id in all syncer

	source            string   // source address
	sourcePassword    string   // source password
	target            []string // target address
	targetPassword    string   // target password
	runId             string   // source runId
	slotLeftBoundary  int      // mark the left slot boundary if enable resuming from break point and is cluster
	slotRightBoundary int      // mark the right slot boundary if enable resuming from break point and is cluster
	httpProfilePort   int      // http profile port

	// stat info
	stat Status

	startDbId                  int    // use in break resume from break-point
	enableResumeFromBreakPoint bool   // enable?
	checkpointName             string // checkpoint name, if is shard, this name has suffix

	/*
	 * this channel is used to calculate delay between redis-shake and target redis.
	 * Once oplog sent, the corresponding delayNode push back into this queue. Next time
	 * receive reply from target redis, the front node poped and then delay calculated.
	 */
	delayChannel chan *delayNode

	sendBuf              chan cmdDetail // sending queue
	WaitFull             chan struct{}  // wait full sync done
	Restart              chan struct{}
	MaxParallelFullSyncs *semaphore.Weighted
	fullSyncRetryCounter int
	lastRetry            time.Time
	sourceOffset         int64
	lastCommittedOffset  int64
}

func (ds *DbSyncer) GetExtraInfo() map[string]interface{} {
	return map[string]interface{}{
		"SourceAddress":      ds.source,
		"TargetAddress":      ds.target,
		"SenderBufCount":     len(ds.sendBuf),
		"ProcessingCmdCount": len(ds.delayChannel),
		"TargetDBOffset":     ds.stat.targetOffset.Get(),
		"SourceDBOffset":     ds.stat.sourceOffset,
	}
}

func (ds *DbSyncer) incrementRetryCounter() {
	if time.Now().After(ds.lastRetry.Add(time.Hour)) {
		ds.fullSyncRetryCounter = 0
	}
	ds.lastRetry = time.Now()
	ds.fullSyncRetryCounter++
	metric.GetMetric(ds.id).RetryCounter = ds.fullSyncRetryCounter
	if ds.fullSyncRetryCounter > 3 {
		log.Panicf("DbSyncer[%d] max amount of failures reached for %s", ds.id, strings.Join(ds.target, ","))
	}
}

// main
func (ds *DbSyncer) Sync() {
	ds.incrementRetryCounter()
	//We bring back the source master offset to the last committed offset to target
	//this is needed as we may loose the internal buffer when Sync() or child calls panic
	ds.sourceOffset = ds.lastCommittedOffset

	log.Infof("DbSyncer[%d] starts syncing data from %v to %v with http[%v], enableResumeFromBreakPoint[%v], "+
		"slot boundary[%v, %v]", ds.id, ds.source, ds.target, ds.httpProfilePort, ds.enableResumeFromBreakPoint,
		ds.slotLeftBoundary, ds.slotRightBoundary)

	var err error
	runId, dbid := "?", 0
	if ds.enableResumeFromBreakPoint {
		// assign the checkpoint name with suffix if is cluster
		if ds.slotLeftBoundary != -1 {
			ds.checkpointName = utils.ChoseSlotInRange(utils.CheckpointKey, ds.slotLeftBoundary, ds.slotRightBoundary)
		}

		// checkpoint reload if has
		log.Infof("DbSyncer[%d] enable resume from break point, try to load checkpoint", ds.id)
		runId, ds.sourceOffset, dbid, err = checkpoint.LoadCheckpoint(ds.id, ds.source, ds.target, conf.Options.TargetAuthType,
			ds.targetPassword, ds.checkpointName, conf.Options.TargetType == conf.RedisTypeCluster, conf.Options.SourceTLSEnable)
		if err != nil {
			log.Panicf("DbSyncer[%d] load checkpoint from %v failed[%v]", ds.id, ds.target, err)
			return
		}
		log.Infof("DbSyncer[%d] checkpoint info: runId[%v], sourceOffset[%v] dbid[%v]", ds.id, runId, ds.sourceOffset, dbid)
	}

	base.Status = "waitfull"
	var input io.ReadCloser
	var nsize int64
	var isFullSync bool

	ds.MaxParallelFullSyncs.Acquire(context.TODO(), 1)

	input, nsize, isFullSync, runId, err = ds.sendPSyncCmd(ds.source, conf.Options.SourceAuthType, ds.sourcePassword,
		conf.Options.SourceTLSEnable, runId)
	if err != nil {
		log.Errorf("DbSyncer[%d] Restarting as error occurred initializing PSYNC %v", err)
		ds.MaxParallelFullSyncs.Release(1)
		go ds.Sync()
		return
	}
	ds.runId = runId

	defer input.Close()

	reader := bufio.NewReaderSize(input, utils.ReaderBufferSize)

	if isFullSync {
		// sync rdb
		log.Infof("DbSyncer[%d] rdb file size = %d", ds.id, nsize)
		base.Status = "full"
		err := ds.syncRDBFile(reader, ds.target, conf.Options.TargetAuthType, ds.targetPassword, nsize, conf.Options.TargetTLSEnable)
		if err != nil {
			log.Errorf("Restarting DbSyncer[%d] for: %v", ds.id, err)
			ds.MaxParallelFullSyncs.Release(1)
			go ds.Sync()
			return
		}
	}
	ds.startDbId = dbid

	log.Infof("DbSyncer[%d] run incr-sync directly with db_id[%v]", ds.id, dbid)
	// set fullSyncProgress to 100 when skip full sync stage
	metric.GetMetric(ds.id).SetFullSyncProgress(ds.id, 100)

	ds.MaxParallelFullSyncs.Release(1)

	// sync increment
	base.Status = "incr"
	close(ds.WaitFull)
	ds.syncCommand(reader, ds.target, conf.Options.TargetAuthType, ds.targetPassword, conf.Options.TargetTLSEnable, dbid)
}
