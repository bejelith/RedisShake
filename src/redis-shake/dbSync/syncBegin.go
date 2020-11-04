package dbSync

import (
	"bufio"
	"github.com/alibaba/RedisShake/pkg/libs/atomic2"
	"github.com/alibaba/RedisShake/pkg/libs/io/pipe"
	"github.com/alibaba/RedisShake/pkg/libs/log"
	"github.com/alibaba/RedisShake/redis-shake/base"
	"github.com/alibaba/RedisShake/redis-shake/common"
	"github.com/pkg/errors"
	"io"
	"net"
	"time"

	"github.com/alibaba/RedisShake/redis-shake/configure"
)

// send command to Source redis

func (ds *DbSyncer) sendSyncCmd(master, authType, passwd string, tlsEnable bool) (net.Conn, int64) {
	c, wait := utils.OpenSyncConn(master, authType, passwd, tlsEnable)
	for {
		select {
		case nsize := <-wait:
			if nsize == 0 {
				log.Infof("DbSyncer[%d] + waiting Source rdb", ds.id)
			} else {
				return c, nsize
			}
		case <-time.After(time.Second):
			log.Infof("DbSyncer[%d] - waiting Source rdb", ds.id)
		}
	}
}

// This function actually doesnt not send Psync command, the command is sent by "SendPSyncContinue" (whatever it means)
func (ds *DbSyncer) sendPSyncCmd(master, authType, passwd string, tlsEnable bool, runId string) (pipe.Reader, int64, bool, string, error) {
	c, err := utils.OpenNetConn(master, authType, passwd, tlsEnable)
	if err != nil {
		return nil, -1, false, "", err
	}
	log.Infof("DbSyncer[%d] psync connect '%v' with auth type[%v] OK!", ds.id, master, authType)

	utils.SendPSyncListeningPort(c, conf.Options.HttpProfile)
	log.Infof("DbSyncer[%d] psync send listening port[%v] OK!", ds.id, conf.Options.HttpProfile)

	// reader buffer bind to client
	br := bufio.NewReaderSize(c, utils.ReaderBufferSize)
	// writer buffer bind to client
	bw := bufio.NewWriterSize(c, utils.WriterBufferSize)

	log.Infof("DbSyncer[%d] try to send 'psync' command: run-Id[%v], sourceOffset[%v]", ds.id, runId, ds.sourceOffset)
	// send psync command and decode the result
	runid, offset, wait, err := utils.SendPSyncContinue(br, bw, runId, ds.sourceOffset)
	if err != nil {
		return nil, 0, false, runid, errors.Errorf("DbSyncer[%d] error sending PSYNC: %v", err)
	}
	ds.sourceOffset = offset
	ds.stat.targetOffset.Set(ds.sourceOffset)

	piper, pipew := pipe.NewSize(utils.ReaderBufferSize)
	if wait == nil {
		// continue
		log.Infof("DbSyncer[%d] psync runid = %s, sourceOffset = %d, psync continue", ds.id, runId, offset)
		go ds.runIncrementalSync(c, br, bw, 0, runid, master, authType, passwd, tlsEnable, pipew, true)
		return piper, 0, false, runid, nil
	} else {
		// fullresync
		log.Infof("DbSyncer[%d] psync runid = %s, sourceOffset = %d, fullsync", ds.id, runid, offset)

		// get rdb file size, wait Source rdb dump successfully.
		var nsize int64
		for nsize == 0 {
			select {
			case nsize = <-wait:
			case <-time.After(time.Second):
			}
			log.Infof("DbSyncer[%d] Waiting fullsync to finish", ds.id)
		}

		go ds.runIncrementalSync(c, br, bw, int(nsize), runid, master, authType, passwd, tlsEnable, pipew, true)
		return piper, nsize, true, runid, nil
	}
}

func (ds *DbSyncer) runIncrementalSync(c net.Conn, br *bufio.Reader, bw *bufio.Writer, rdbSize int, runId string,
	master, authType, passwd string, tlsEnable bool, pipew pipe.Writer, isFullSync bool) {
	// write -> pipew -> piper -> read
	defer pipew.Close()
	if isFullSync { //always true -.-
		log.Info("DbSyncer[%d] RDB size %d", ds.id, rdbSize)
		p := make([]byte, 8192)
		// read rdb in for loop
		for rdbSize != 0 {
			// br -> pipew
			rdbSize -= utils.Iocopy(br, pipew, p, rdbSize)
		}
	}

	for {
		/*
		 * read from br(Source redis) and write into pipew.
		 * Generally speaking, this function is forever run.
		 */
		_, err := ds.pSyncPipeCopy(c, br, bw, pipew)
		if err != nil {
			// Do not exists c (Connection) will be rebuild later in the loop
			log.Errorf("DbSyncer[%d] psync runid = %s, sourceOffset = %d, pipe is broken",
				ds.id, runId, ds.sourceOffset, err)
		}
		// We dont leave this for loop anymore when an error occurs so we need to increase the retry counter
		ds.incrementRetryCounter()
		// the 'c' is closed every loop

		ds.stat.targetOffset.Set(ds.sourceOffset)

		// reopen 'c' every time
		for {
			// ds.SyncStat.SetStatus("reopen")
			base.Status = "reopen"
			time.Sleep(time.Second)

			c = utils.OpenNetConnSoft(master, authType, passwd, tlsEnable) //TODO we dont get to know what's the error when open a socket?!
			if c != nil {
				log.Infof("DbSyncer[%d] Event:SourceConnReopenSuccess\tId: %s\tsourceOffset = %d",
					ds.id, conf.Options.Id, ds.sourceOffset)
				base.Status = "incr"
				utils.AuthPassword(c, authType, passwd)
				utils.SendPSyncListeningPort(c, conf.Options.HttpProfile)
				br = bufio.NewReaderSize(c, utils.ReaderBufferSize)
				bw = bufio.NewWriterSize(c, utils.WriterBufferSize)
				_, _, _, err = utils.SendPSyncContinue(br, bw, runId, ds.sourceOffset)
				if err != nil {
					// If PSYNC fails we need skip PipeCopy but retry to re-establish connection to the Source so we stay in the loop
					log.Errorf("DbSyncer[%d] retrying incremental sync as was error found %d", err)
					time.Sleep(30 * time.Second) //TODO maybe implement an exponential backoff
				}

				//Break the loop to start PipeCopy again
				break
			} else {
				log.Errorf("DbSyncer[%d] Event:SourceConnReopenFail\tId: %s", ds.id, conf.Options.Id)
			}
		}

	}
}

func (ds *DbSyncer) pSyncPipeCopy(c net.Conn, br *bufio.Reader, bw *bufio.Writer, copyto io.Writer) (int64, error) {
	var nread atomic2.Int64
	go func() {
		defer c.Close()
		for range time.NewTicker(time.Second).C {
			select {
			case <-ds.WaitFull:
				ds.sourceOffset += nread.Get()
				if err := utils.SendPSyncAck(bw, ds.sourceOffset); err != nil {
					log.Errorf("dbSyncer[%v] send sourceOffset to Source redis failed[%v]", ds.id, err)
					return
				}
			default:
				if err := utils.SendPSyncAck(bw, 0); err != nil {
					log.Errorf("dbSyncer[%v] send sourceOffset to Source redis failed[%v]", ds.id, err)
					return
				}
			}
		}
	}()

	var p = make([]byte, 8192)
	for {
		n, err := br.Read(p)
		if err != nil {
			return nread.Get(), err
		}
		if _, err := copyto.Write(p[:n]); err != nil {
			return nread.Get(), err
		}
		nread.Add(int64(n))
	}
}
