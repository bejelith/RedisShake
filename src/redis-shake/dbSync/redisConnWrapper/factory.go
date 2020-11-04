package redisConnWrapper

import (
	utils "github.com/alibaba/RedisShake/redis-shake/common"
	redigo "github.com/garyburd/redigo/redis"
)

type RedisConnFactory func(host, password string, tlsEnable bool) (redigo.Conn, error)

func DefaultRedisConnFactory(host, password string, tlsEnable bool) (redigo.Conn, error) {
	conn, err := utils.OpenNetConn(host, "auth", password, tlsEnable)
	if err != nil {
		return nil, err
	}
	return redigo.NewConn(conn, 0, 0), nil
}

type RedisClusterFactory func(masters []string, password string, tlsEnable bool) (redigo.Conn, error)

func DefaultRedisClusterFactory(masters []string, password string, tlsEnable bool) (redigo.Conn, error) {
	conn, err := utils.OpenRedisConn(masters, "auth", password, true, tlsEnable)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
