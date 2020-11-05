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

type RedisClusterFactory func(masters []string, password string, tlsEnable bool) (*redis.Cluster, error)

//func DefaultRedisClusterFactory(masters []string, password string, tlsEnable bool) (redigo.Conn, error) {
//	conn, err := utils.OpenRedisConn(masters, "auth", password, true, tlsEnable)
//	if err != nil {
//		return nil, err
//	}
//	return conn, nil
//}

func DefaultRedisClusterFactory(masters []string, password string, tlsEnable bool) (*redis.Cluster, error) {
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
