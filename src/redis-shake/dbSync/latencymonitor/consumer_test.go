package latencymonitor

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"

	"github.com/alibaba/RedisShake/redis-shake/metric"
)

func TestKeyFilterRegexp(t *testing.T) {
	cases := map[string]bool{
		"1":    true,
		"1a":   false,
		"":     false,
		"1110": true,
	}
	for k, v := range cases {
		assert.Equal(t, v, keyPrefixRegex.MatchString(keyPrefix+k), k, keyPrefixRegex.String())
	}
}

func TestCalculateLatench(t *testing.T) {
	metricId := 1
	metric.AddMetric(metricId)
	cases := [][]string{
		{"set", "key1", "123", "0"},
		{"set", keyPrefix + "1", strconv.Itoa(int(time.Now().Add(-time.Second).UnixNano())), "1000000000"},
	}
	for _, args := range cases {
		stringArgs := make([][]byte, 0)
		for _, arg := range args {
			stringArgs = append(stringArgs, []byte(arg))
		}
		CalcLatency(args[0], stringArgs[1:3], metricId)
		expected, _ := strconv.Atoi(string(stringArgs[3]))
		actual := int(metric.GetMetric(metricId).SourceLatency.Get())
		if expected != 0 {
			assert.Greater(t, actual, expected, "Key: %s", stringArgs[1])
		} else {
			assert.Zero(t, expected)
		}
		metric.GetMetric(metricId).SourceLatency.Set(0)
	}
}
