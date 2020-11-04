package latencymonitor

import (
	"fmt"
	"github.com/alibaba/RedisShake/pkg/libs/atomic2"
	"github.com/alibaba/RedisShake/pkg/libs/log"
	utils "github.com/alibaba/RedisShake/redis-shake/common"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/redisConnWrapper"
<<<<<<< HEAD
	"regexp"
=======
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
	"strconv"
	"time"
)

<<<<<<< HEAD
var keyPrefix = "synthetic_latency_generator_"
var keyPrefixRegex, _ = regexp.Compile("^"+ keyPrefix +"\\d+$")
var tickDuration = 15 * time.Second
=======
var KEY_PREFIX = "sinthetic_latency_generator_"
>>>>>>> 6643c58 (Produce keys for each master with timestamp)

type Producer interface {
	Run()
	Stop()
	Error() error
}

func NewSyntheticProducer(slots []utils.SlotOwner, password string, tls bool) Producer {
	keys, masters := calculateKeys(slots)
	return &producer{
<<<<<<< HEAD
		slots:                     slots,
		keys:                      keys,
		masters:                   masters,
		password:                  password,
		tls:                       tls,
		runChannel:                make(chan struct{}),
		running:                   atomic2.Bool{},
=======
		slots:      slots,
		keys:       keys,
		masters:    masters,
		password:   password,
		tls:        tls,
		runChannel: make(chan struct{}),
		running:    atomic2.Bool{},
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
		redisClusterClientFactory: redisConnWrapper.DefaultRedisClusterFactory,
	}
}

type producer struct {
<<<<<<< HEAD
	slots                     []utils.SlotOwner
	keys                      []string
	runChannel                chan struct{}
	password                  string
	tls                       bool
	masters                   []string
	error                     error
	running                   atomic2.Bool
	redisClusterClientFactory redisConnWrapper.RedisClusterFactory
	client                    redisConnWrapper.ClusterI
=======
	slots      []utils.SlotOwner
	keys       []string
	runChannel chan struct{}
	password   string
	tls        bool
	masters    []string
	error      error
	running    atomic2.Bool
	redisClusterClientFactory redisConnWrapper.RedisClusterFactory
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
}

func (p *producer) Error() error {
	return p.error
}

func findKeyInRange(min, max int) string {
	for i := 0; ; i++ {
<<<<<<< HEAD
		key := fmt.Sprintf("%s%s", keyPrefix, strconv.Itoa(i))
		hash := int(crc16(key) & (16384 - 1))
=======
		key := fmt.Sprintf("%s%s", KEY_PREFIX, strconv.Itoa(i))
		hash := int(crc16(key))
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
		if hash >= min && hash <= max {
			return key
		}
	}
}

func (p *producer) run() {
<<<<<<< HEAD
	defer p.Stop()
	var err error
	p.client, err = p.redisClusterClientFactory(p.masters, p.password, p.tls)
	if err != nil {
		log.Errorf("SyntheticProducer interrupted for error %v", err)
		p.error = err
		return
	}
	ticker := time.NewTicker(tickDuration)
	defer p.client.Close()
	for {
		select {
		case <-ticker.C:
			for _, key := range p.keys {
				now := strconv.Itoa(int(time.Now().UnixNano()))
				if _, err := p.client.Do("set", key, now, "EX", strconv.Itoa(int(tickDuration/time.Second))); err != nil {
					log.Warnf("SyntheticProducer failed to update key %s for %v", key, err)
				}
			}
			//if err := c.Flush(); err != nil {
			//	log.Warnf("SyntheticProducer error found when flushing commands to target cluster %v", err)
			//}
		case <-p.runChannel:
			log.Info("SyntheticProducer stopping")
			return
		}
=======
	c, err := p.redisClusterClientFactory(p.masters, p.password, p.tls)
	if err != nil {
		log.Errorf("Synthetic producer stopped for %v", err)
		p.error = err
		return
	}
	ticker := time.NewTicker(15 * time.Second)
	defer c.Close()
	defer p.Stop()
	select {
	case <- ticker.C:
		for key := range p.keys {
			now := strconv.Itoa(int(time.Now().UnixNano()))
			if _, err := c.Do("set", key, now); err != nil {
				log.Warn("Latency monitor failed to update key %s for %v", key, err)
			}
		}
	case <- p.runChannel:
		break
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
	}
}

func (p *producer) Run() {
	if !p.running.CompareAndSwap(false, true) {
		return
	}
<<<<<<< HEAD
	log.Infof("SyntheticProducer started producing timestamps for source cluster %v with keys %v", p.masters, p.keys)
=======
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
	p.runChannel = make(chan struct{})
	go p.run()
}

func (p *producer) Stop() {
<<<<<<< HEAD
	log.Info("SyntheticProducer stop called, waiting to finish")
	if p.running.CompareAndSwap(true, false) {
		close(p.runChannel)
	}
=======
	close(p.runChannel)
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
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
