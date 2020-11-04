package redisConnWrapper

type MockRedisConn struct {
	DoCount int
	Host    string
	DoFunc      func(string, ...interface{}) (interface{}, error)
}

func (m MockRedisConn) Close() error {
	return nil
}

func (m MockRedisConn) Err() error {
	panic("implement me")
}

func (m MockRedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return m.DoFunc(commandName, args...)
}

func (m MockRedisConn) Send(_ string, _ ...interface{}) error {
	panic("implement me")
}

func (m MockRedisConn) Flush() error {
	panic("implement me")
}

func (m MockRedisConn) Receive() (reply interface{}, err error) {
	panic("implement me")
}

