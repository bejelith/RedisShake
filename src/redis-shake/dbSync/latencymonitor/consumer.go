package latencymonitor

import (
	"github.com/alibaba/RedisShake/pkg/libs/log"
	"github.com/alibaba/RedisShake/redis-shake/metric"
	"strconv"
	"time"
)

func CalcLatency(cmd string, args [][]byte, dsId int) {
	if cmd != "set" || len(args) < 2 {
		return
	}
	key := string(args[0])
	if keyPrefixRegex.MatchString(key) {
		value := string(args[1])
		messageTime, err := strconv.Atoi(value)
		if err != nil {
			log.Warnf("CalcLatency: Invalid latency value read from key %s: %v", key, err)
			return
		}
		delta := time.Now().UnixNano() - int64(messageTime)
		metric.GetMetric(dsId).SourceLatency.Add(delta)
	}
}
