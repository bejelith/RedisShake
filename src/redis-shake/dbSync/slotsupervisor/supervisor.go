package slotsupervisor

import (
	"github.com/alibaba/RedisShake/pkg/libs/errors"
	"github.com/alibaba/RedisShake/pkg/libs/log"
	conf "github.com/alibaba/RedisShake/redis-shake/configure"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/redisConnWrapper"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/slot"
	redigo "github.com/garyburd/redigo/redis"
	"regexp"
	"strings"
	"time"
)

type SlotSupervisor interface {
	GetSlotState() (*slot.SyncNode, error)
}

type slotSupervisor struct {
	slot             slot.SyncNode
	redisConnFactory redisConnWrapper.RedisConnFactory
	maxRetries       int
}

func New(slot slot.SyncNode) SlotSupervisor {
	ss := &slotSupervisor{
		slot:             slot,
		redisConnFactory: redisConnWrapper.DefaultRedisConnFactory,
		maxRetries:       6,
	}
	return ss
}

func (s *slotSupervisor) GetSlotState() (*slot.SyncNode, error) {
	return s.recursiveGetSlotState(s.maxRetries)
}

func (s *slotSupervisor) recursiveGetSlotState(recursionDept int) (*slot.SyncNode, error) {
	newSlot := s.slot           //Copy slow info
	newSlot.Slaves = []string{} // Reset slaves
	masterFound := false
	var err error
	hosts := append([]string{s.slot.Source}, s.slot.Slaves...)
	for _, host := range hosts {
		isMaster := false
		isMaster, err = s.getRedisNodeState(host, s.slot.SourcePassword, conf.Options.SourceTLSEnable)
		if isMaster && err == nil {
			masterFound = isMaster
			newSlot.Source = host
		} else {
			newSlot.Slaves = append(newSlot.Slaves, host)
			if err != nil {
				log.Errorf("GetSlotState - error while discovering slot topology: %v", err)
			}
		}
	}
	if masterFound {
		log.Infof("Master node %s found for slot %d/%d", newSlot.Source, newSlot.SlotLeftBoundary, newSlot.SlotRightBoundary)
		return &newSlot, nil
	} else if recursionDept == 0 {
		return nil, errors.New("Max retries reached")
	}
	time.Sleep(time.Second * time.Duration(recursionDept))
	return s.recursiveGetSlotState(recursionDept - 1)
}

func (s *slotSupervisor) getRedisNodeState(node, password string, tlsEnable bool) (bool, error) {
	masterRegex := regexp.MustCompile("^role:master")
	slaveRegex := regexp.MustCompile("^role:slave")
	conn, connErr := s.redisConnFactory(node, password, tlsEnable)
	if connErr != nil {
		return false, connErr
	}
	defer conn.Close()
	resp, err := redigo.String(conn.Do("info", "replication"))
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(resp, "\n") {
		log.Infof("%s Testing line %s", node, line)
		if masterRegex.MatchString(line) {
			return true, nil
		} else if slaveRegex.MatchString(line) {
			return false, nil
		}
	}
	return false, errors.Errorf("Redis node %s is reporting a invalid role", node)
}
