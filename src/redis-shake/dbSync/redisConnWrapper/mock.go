package redisConnWrapper

type MockRedisConn struct {
	DoCount int
	Host    string
<<<<<<< HEAD
	DoFunc  func(string, ...interface{}) (interface{}, error)
=======
	DoFunc      func(string, ...interface{}) (interface{}, error)
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
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

<<<<<<< HEAD
type MockRedisCluster struct {
	DoCount int
	Host    string
	DoFunc  func(string, ...interface{}) (interface{}, error)
}

func (m MockRedisCluster) Close() {
	return
}

func (m MockRedisCluster) Err() error {
	panic("implement me")
}

func (m MockRedisCluster) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	m.DoCount = m.DoCount + 1
	return m.DoFunc(commandName, args...)
}

func (m MockRedisCluster) Send(_ string, _ ...interface{}) error {
	panic("implement me")
}

func (m MockRedisCluster) Flush() error {
	panic("implement me")
}

func (m MockRedisCluster) Receive() (reply interface{}, err error) {
	panic("implement me")
}
=======
>>>>>>> 6643c58 (Produce keys for each master with timestamp)
