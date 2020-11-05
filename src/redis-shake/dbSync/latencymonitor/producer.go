package latencymonitor

import (
	"fmt"
	"github.com/alibaba/RedisShake/pkg/libs/atomic2"
	"github.com/alibaba/RedisShake/pkg/libs/log"
	utils "github.com/alibaba/RedisShake/redis-shake/common"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/redisConnWrapper"
	"strconv"
	"time"
)

var KeyPrefix = "synthetic_latency_generator_"

type Producer interface {
	Run()
	Stop()
	Error() error
}

func NewSyntheticProducer(slots []utils.SlotOwner, password string, tls bool) Producer {
	keys, masters := calculateKeys(slots)
	return &producer{
		slots:                     slots,
		keys:                      keys,
		masters:                   masters,
		password:                  password,
		tls:                       tls,
		runChannel:                make(chan struct{}),
		running:                   atomic2.Bool{},
		redisClusterClientFactory: redisConnWrapper.DefaultRedisClusterFactory,
	}
}

type producer struct {
	slots                     []utils.SlotOwner
	keys                      []string
	runChannel                chan struct{}
	password                  string
	tls                       bool
	masters                   []string
	error                     error
	running                   atomic2.Bool
	redisClusterClientFactory redisConnWrapper.RedisClusterFactory
}

func (p *producer) Error() error {
	return p.error
}

func findKeyInRange(min, max int) string {
	for i := 0; ; i++ {
		key := fmt.Sprintf("%s%s", KeyPrefix, strconv.Itoa(i))
		hash := int(crc16(key))
		if hash >= min && hash <= max {
			return key
		}
	}
}

func (p *producer) run() {
	defer p.Stop()
	c, err := p.redisClusterClientFactory(p.masters, p.password, p.tls)
	if err != nil {
		log.Errorf("SyntheticProducer interrupted for error %v", err)
		p.error = err
		return
	}
	ticker := time.NewTicker(15 * time.Second)
	defer c.Close()
	for {
		select {
		case <-ticker.C:
			for key := range p.keys {
				now := strconv.Itoa(int(time.Now().UnixNano()))
				if _, err := c.Do("set", key, now); err != nil {
					log.Warn("SyntheticProducer failed to update key %s for %v", key, err)
				}
			}
			c.Flush()
		case <-p.runChannel:
			log.Info("SyntheticProducer stopping")
			break
		}
	}
}

func (p *producer) Run() {
	if !p.running.CompareAndSwap(false, true) {
		return
	}
	log.Info("SyntheticProducer started producing timestamps")
	p.runChannel = make(chan struct{})
	go p.run()
}

func (p *producer) Stop() {
	log.Info("SyntheticProducer stop called, waiting to finish")
	p.running.Set(false)
	close(p.runChannel)
}

func calculateKeys(slots []utils.SlotOwner) ([]string, []string) {
	var keys []string
	var masters []string
	for _, slot := range slots {
		keys = append(keys, findKeyInRange(slot.SlotLeftBoundary, slot.SlotRightBoundary))
		masters = append(masters, slot.Master)
	}
	return keys, masters
}
