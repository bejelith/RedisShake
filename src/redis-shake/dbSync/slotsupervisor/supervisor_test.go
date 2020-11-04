package slotsupervisor

import (
	"errors"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/redisConnWrapper"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/slot"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindAValidMaster(t *testing.T) {
	redisConnFactory := func(host, passwd string, tls bool) (redigo.Conn, error) {
		mock := redisConnWrapper.MockRedisConn{
			Host: host,
			DoFunc: func(commandName string, args ...interface{}) (interface{}, error) {
				if host == "master" {
					return true, nil
				}
				return nil, errors.New("Mock Error")
			},
		}
		return mock, nil
	}
	baseSlot := slot.SyncNode{
		Source:            "master",
		Slaves:            []string{"slave1", "slave2"},
		SourcePassword:    "pwd",
		SlotLeftBoundary:  0,
		SlotRightBoundary: 100,
	}
	supervisor := New(baseSlot).(*slotSupervisor)
	supervisor.redisConnFactory = redisConnFactory
	newSlot, err := supervisor.GetSlotState()
	assert.Nil(t, err, "Error found!")
	assert.Equal(t, "master", newSlot.Source)
}

func TestNoAvailableMasterFound(t *testing.T) {
	redisConnFactory := func(a, b string, c bool) (redigo.Conn, error) {
		mockConn := redisConnWrapper.MockRedisConn{
			Host: a,
			DoFunc: func(commandName string, args ...interface{}) (interface{}, error) {
				return nil, errors.New("Mock Error")
			},
		}
		return mockConn, nil
	}
	baseSlot := slot.SyncNode{
		Source:            "master",
		Slaves:            []string{"slave1", "slave2"},
		SourcePassword:    "pwd",
		SlotLeftBoundary:  0,
		SlotRightBoundary: 100,
	}
	supervisor := New(baseSlot).(*slotSupervisor)
	supervisor.redisConnFactory = redisConnFactory
	newSlot, err := supervisor.GetSlotState()
	assert.NotNilf(t, err, "Error should not be nil")
	assert.Nil(t, newSlot, "Result should be nil")
}
