package redisConnWrapper

import (
	utils "github.com/alibaba/RedisShake/redis-shake/common"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/vinllen/redis-go-cluster"
	"time"
)

type RedisConnFactory func(host, password string, tlsEnable bool) (redigo.Conn, error)

func DefaultRedisConnFactory(host, password string, tlsEnable bool) (redigo.Conn, error) {
	conn, err := utils.OpenNetConn(host, "auth", password, tlsEnable)
	if err != nil {
		return nil, err
	}
	return redigo.NewConn(conn, 0, 0), nil
}

type ClusterI interface{
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
	Close()
}

type RedisClusterFactory func(masters []string, password string, tlsEnable bool) (ClusterI, error)

func DefaultRedisClusterFactory(masters []string, password string, tlsEnable bool) (ClusterI, error) {
	options := &redis.Options{
		StartNodes:  masters,
		ConnTimeout: time.Second * 2,
		ReadTimeout: time.Second / 2,
		Password:    password,
	}
	conn, err := redis.NewCluster(options)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
