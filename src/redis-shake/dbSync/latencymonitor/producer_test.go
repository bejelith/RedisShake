package latencymonitor

import (
	"github.com/alibaba/RedisShake/redis-shake/dbSync/redisConnWrapper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func mockFactory(mock redisConnWrapper.MockRedisCluster) redisConnWrapper.RedisClusterFactory {
	return func(masters []string, password string, tlsEnable bool) (redisConnWrapper.ClusterI, error) {
		return &mock, nil
	}
}

func TestClusterClientIsCalled(t *testing.T) {
	var doCount = 0
	var doArgs []interface{}
	tickDuration = time.Second
	mock := redisConnWrapper.MockRedisCluster{
		DoFunc: func(s string, args ...interface{}) (interface{}, error) {
			doCount = doCount + 1
			doArgs = args
			return nil, nil
		},
	}
	p := producer{
		keys: []string{"key"},
		redisClusterClientFactory: mockFactory(mock),
	}
	p.Run()
	time.Sleep(tickDuration+tickDuration/2)
	p.Stop()
	assert.Equal(t, 1, doCount)
	assert.Equal(t, "key", doArgs[0])
	assert.Equal(t, "EX", doArgs[2])
	assert.Equal(t, "1", doArgs[3])
	assert.False(t, p.running.Get())
}
