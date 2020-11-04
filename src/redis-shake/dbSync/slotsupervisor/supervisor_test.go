package slotsupervisor

import (
	"errors"
	"github.com/alibaba/RedisShake/redis-shake/dbSync/slot"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestR(t *testing.T){
	slaveRegex := regexp.MustCompile("^role:slave")
	line := "role:slave \r"
	assert.True(t, slaveRegex.MatchString(line))
}

func TestFindAValidMaster(t *testing.T) {
	defaultRedisConnFactory = func(host, passwd string, tls bool) (redigo.Conn, error) {
		mock := mockRedisConn{
			host: host,
			do: func(commandName string, args ...interface{}) (interface{}, error) {
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
	supervisor := New(baseSlot)
	newSlot, err := supervisor.GetSlotState()
	assert.Nil(t, err, "Error found!")
	assert.Equal(t, "master", newSlot.Source)
}

func TestNoAvailableMasterFound(t *testing.T) {
	defaultRedisConnFactory = func(a, b string, c bool) (redigo.Conn, error) {
		mockConn := mockRedisConn{
			host: a,
			do: func(commandName string, args ...interface{}) (interface{}, error) {
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
	supervisor := New(baseSlot)
	newSlot, err := supervisor.GetSlotState()
	assert.NotNilf(t, err, "Error should not be nil")
	assert.Nil(t, newSlot, "Result should be nil")
}

type mockRedisConn struct {
	doCount int
	host    string
	do      func(string, ...interface{}) (interface{}, error)
}

func (m mockRedisConn) Close() error {
	panic("implement me")
}

func (m mockRedisConn) Err() error {
	panic("implement me")
}

func (m mockRedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return m.do(commandName, args...)
}

func (m mockRedisConn) Send(_ string, _ ...interface{}) error {
	panic("implement me")
}

func (m mockRedisConn) Flush() error {
	panic("implement me")
}

func (m mockRedisConn) Receive() (reply interface{}, err error) {
	panic("implement me")
}
