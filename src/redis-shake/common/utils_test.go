// +build integration
// +build linux darwin windows

package utils

import (
	"fmt"
	"testing"

	"github.com/alibaba/RedisShake/redis-shake/unit_test_common"

	"github.com/stretchr/testify/assert"
)

var (
	testAddrCluster = unit_test_common.TestUrlCluster
)

/*func TestGetAllClusterNode(t *testing.T) {
	var nr int
	{
		fmt.Printf("TestGetAllClusterNode case %d.\n", nr)
		nr++

		client := OpenRedisConn([]string{"10.1.1.1:21333"}, "auth", "123456", false, false)
		ret, err := GetAllClusterNode(client, "master", "")
		sort.Strings(ret)
		assert.Equal(t, nil, err, "should be equal")
		assert.Equal(t, 3, len(ret), "should be equal")
		assert.Equal(t, "10.1.1.1:21331", ret[0], "should be equal")
		assert.Equal(t, "10.1.1.1:21332", ret[1], "should be equal")
		assert.Equal(t, "10.1.1.1:21333", ret[2], "should be equal")
	}

	{
		fmt.Printf("TestGetAllClusterNode case %d.\n", nr)
		nr++

		client := OpenRedisConn([]string{"10.1.1.1:21333"}, "auth", "123456", false, false)
		ret, err := GetAllClusterNode(client, "slave", "")
		sort.Strings(ret)
		assert.Equal(t, nil, err, "should be equal")
		assert.Equal(t, 3, len(ret), "should be equal")
		assert.Equal(t, "10.1.1.1:21334", ret[0], "should be equal")
		assert.Equal(t, "10.1.1.1:21335", ret[1], "should be equal")
		assert.Equal(t, "10.1.1.1:21336", ret[2], "should be equal")
	}

	{
		fmt.Printf("TestGetAllClusterNode case %d.\n", nr)
		nr++

		client := OpenRedisConn([]string{"10.1.1.1:21333"}, "auth", "123456", false, false)
		ret, err := GetAllClusterNode(client, "all", "")
		sort.Strings(ret)
		assert.Equal(t, nil, err, "should be equal")
		assert.Equal(t, 6, len(ret), "should be equal")
		assert.Equal(t, "10.1.1.1:21331", ret[0], "should be equal")
		assert.Equal(t, "10.1.1.1:21332", ret[1], "should be equal")
		assert.Equal(t, "10.1.1.1:21333", ret[2], "should be equal")
		assert.Equal(t, "10.1.1.1:21334", ret[3], "should be equal")
		assert.Equal(t, "10.1.1.1:21335", ret[4], "should be equal")
		assert.Equal(t, "10.1.1.1:21336", ret[5], "should be equal")
	}
}*/

func TestCutRedisInfoSegment(t *testing.T) {
	input := []byte{35, 32, 83, 101, 114, 118, 101, 114, 13, 10, 114, 101, 100, 105, 115, 95, 118, 101, 114, 115, 105, 111, 110, 58, 51, 46, 50, 46, 49, 49, 13, 10, 114, 101, 100, 105, 115, 95, 103, 105, 116, 95, 115, 104, 97, 49, 58, 100, 101, 49, 97, 100, 48, 50, 54, 13, 10, 114, 101, 100, 105, 115, 95, 103, 105, 116, 95, 100, 105, 114, 116, 121, 58, 48, 13, 10, 114, 101, 100, 105, 115, 95, 98, 117, 105, 108, 100, 95, 105, 100, 58, 54, 101, 49, 99, 101, 51, 55, 97, 102, 101, 48, 54, 99, 51, 52, 99, 13, 10, 114, 101, 100, 105, 115, 95, 109, 111, 100, 101, 58, 115, 116, 97, 110, 100, 97, 108, 111, 110, 101, 13, 10, 111, 115, 58, 76, 105, 110, 117, 120, 32, 51, 46, 49, 48, 46, 48, 45, 57, 53, 55, 46, 50, 49, 46, 51, 46, 101, 108, 55, 46, 120, 56, 54, 95, 54, 52, 32, 120, 56, 54, 95, 54, 52, 13, 10, 97, 114, 99, 104, 95, 98, 105, 116, 115, 58, 54, 52, 13, 10, 109, 117, 108, 116, 105, 112, 108, 101, 120, 105, 110, 103, 95, 97, 112, 105, 58, 101, 112, 111, 108, 108, 13, 10, 103, 99, 99, 95, 118, 101, 114, 115, 105, 111, 110, 58, 52, 46, 56, 46, 53, 13, 10, 112, 114, 111, 99, 101, 115, 115, 95, 105, 100, 58, 57, 48, 54, 57, 13, 10, 114, 117, 110, 95, 105, 100, 58, 49, 52, 49, 98, 51, 54, 101, 102, 50, 51, 50, 57, 56, 52, 53, 54, 98, 52, 100, 97, 101, 48, 49, 51, 49, 50, 57, 49, 98, 50, 100, 101, 99, 99, 101, 102, 102, 52, 54, 53, 13, 10, 116, 99, 112, 95, 112, 111, 114, 116, 58, 56, 53, 48, 51, 13, 10, 117, 112, 116, 105, 109, 101, 95, 105, 110, 95, 115, 101, 99, 111, 110, 100, 115, 58, 53, 56, 57, 54, 51, 49, 13, 10, 117, 112, 116, 105, 109, 101, 95, 105, 110, 95, 100, 97, 121, 115, 58, 54, 13, 10, 104, 122, 58, 49, 48, 13, 10, 108, 114, 117, 95, 99, 108, 111, 99, 107, 58, 51, 55, 52, 50, 48, 50, 48, 13, 10, 101, 120, 101, 99, 117, 116, 97, 98, 108, 101, 58, 47, 100, 97, 116, 97, 47, 109, 97, 112, 103, 111, 111, 47, 99, 111, 100, 105, 115, 47, 46, 47, 98, 105, 110, 47, 99, 111, 100, 105, 115, 45, 115, 101, 114, 118, 101, 114, 13, 10, 99, 111, 110, 102, 105, 103, 95, 102, 105, 108, 101, 58, 47, 100, 97, 116, 97, 47, 109, 97, 112, 103, 111, 111, 47, 99, 111, 100, 105, 115, 47, 46, 47, 99, 111, 100, 105, 115, 95, 115, 101, 114, 118, 101, 114, 47, 99, 111, 110, 102, 47, 114, 101, 100, 105, 115, 45, 56, 53, 48, 51, 46, 99, 111, 110, 102, 13, 10, 13, 10, 35, 32, 67, 108, 105, 101, 110, 116, 115, 13, 10, 99, 111, 110, 110, 101, 99, 116, 101, 100, 95, 99, 108, 105, 101, 110, 116, 115, 58, 52, 57, 13, 10, 99, 108, 105, 101, 110, 116, 95, 108, 111, 110, 103, 101, 115, 116, 95, 111, 117, 116, 112, 117, 116, 95, 108, 105, 115, 116, 58, 48, 13, 10, 99, 108, 105, 101, 110, 116, 95, 98, 105, 103, 103, 101, 115, 116, 95, 105, 110, 112, 117, 116, 95, 98, 117, 102, 58, 48, 13, 10, 98, 108, 111, 99, 107, 101, 100, 95, 99, 108, 105, 101, 110, 116, 115, 58, 48, 13, 10, 13, 10, 35, 32, 77, 101, 109, 111, 114, 121, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 58, 51, 55, 51, 55, 49, 53, 50, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 95, 104, 117, 109, 97, 110, 58, 51, 46, 53, 54, 77, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 95, 114, 115, 115, 58, 52, 57, 51, 57, 55, 55, 54, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 95, 114, 115, 115, 95, 104, 117, 109, 97, 110, 58, 52, 46, 55, 49, 77, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 95, 112, 101, 97, 107, 58, 53, 51, 55, 52, 52, 48, 48, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 95, 112, 101, 97, 107, 95, 104, 117, 109, 97, 110, 58, 53, 46, 49, 51, 77, 13, 10, 116, 111, 116, 97, 108, 95, 115, 121, 115, 116, 101, 109, 95, 109, 101, 109, 111, 114, 121, 58, 54, 55, 51, 56, 55, 50, 56, 53, 53, 48, 52, 13, 10, 116, 111, 116, 97, 108, 95, 115, 121, 115, 116, 101, 109, 95, 109, 101, 109, 111, 114, 121, 95, 104, 117, 109, 97, 110, 58, 54, 50, 46, 55, 54, 71, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 95, 108, 117, 97, 58, 51, 55, 56, 56, 56, 13, 10, 117, 115, 101, 100, 95, 109, 101, 109, 111, 114, 121, 95, 108, 117, 97, 95, 104, 117, 109, 97, 110, 58, 51, 55, 46, 48, 48, 75, 13, 10, 109, 97, 120, 109, 101, 109, 111, 114, 121, 58, 48, 13, 10, 109, 97, 120, 109, 101, 109, 111, 114, 121, 95, 104, 117, 109, 97, 110, 58, 48, 66, 13, 10, 109, 97, 120, 109, 101, 109, 111, 114, 121, 95, 112, 111, 108, 105, 99, 121, 58, 110, 111, 101, 118, 105, 99, 116, 105, 111, 110, 13, 10, 109, 101, 109, 95, 102, 114, 97, 103, 109, 101, 110, 116, 97, 116, 105, 111, 110, 95, 114, 97, 116, 105, 111, 58, 49, 46, 51, 50, 13, 10, 109, 101, 109, 95, 97, 108, 108, 111, 99, 97, 116, 111, 114, 58, 106, 101, 109, 97, 108, 108, 111, 99, 45, 52, 46, 48, 46, 51, 13, 10, 13, 10, 35, 32, 80, 101, 114, 115, 105, 115, 116, 101, 110, 99, 101, 13, 10, 108, 111, 97, 100, 105, 110, 103, 58, 48, 13, 10, 114, 100, 98, 95, 99, 104, 97, 110, 103, 101, 115, 95, 115, 105, 110, 99, 101, 95, 108, 97, 115, 116, 95, 115, 97, 118, 101, 58, 48, 13, 10, 114, 100, 98, 95, 98, 103, 115, 97, 118, 101, 95, 105, 110, 95, 112, 114, 111, 103, 114, 101, 115, 115, 58, 48, 13, 10, 114, 100, 98, 95, 108, 97, 115, 116, 95, 115, 97, 118, 101, 95, 116, 105, 109, 101, 58, 49, 53, 54, 52, 48, 50, 49, 52, 57, 54, 13, 10, 114, 100, 98, 95, 108, 97, 115, 116, 95, 98, 103, 115, 97, 118, 101, 95, 115, 116, 97, 116, 117, 115, 58, 111, 107, 13, 10, 114, 100, 98, 95, 108, 97, 115, 116, 95, 98, 103, 115, 97, 118, 101, 95, 116, 105, 109, 101, 95, 115, 101, 99, 58, 48, 13, 10, 114, 100, 98, 95, 99, 117, 114, 114, 101, 110, 116, 95, 98, 103, 115, 97, 118, 101, 95, 116, 105, 109, 101, 95, 115, 101, 99, 58, 45, 49, 13, 10, 97, 111, 102, 95, 101, 110, 97, 98, 108, 101, 100, 58, 49, 13, 10, 97, 111, 102, 95, 114, 101, 119, 114, 105, 116, 101, 95, 105, 110, 95, 112, 114, 111, 103, 114, 101, 115, 115, 58, 48, 13, 10, 97, 111, 102, 95, 114, 101, 119, 114, 105, 116, 101, 95, 115, 99, 104, 101, 100, 117, 108, 101, 100, 58, 48, 13, 10, 97, 111, 102, 95, 108, 97, 115, 116, 95, 114, 101, 119, 114, 105, 116, 101, 95, 116, 105, 109, 101, 95, 115, 101, 99, 58, 45, 49, 13, 10, 97, 111, 102, 95, 99, 117, 114, 114, 101, 110, 116, 95, 114, 101, 119, 114, 105, 116, 101, 95, 116, 105, 109, 101, 95, 115, 101, 99, 58, 45, 49, 13, 10, 97, 111, 102, 95, 108, 97, 115, 116, 95, 98, 103, 114, 101, 119, 114, 105, 116, 101, 95, 115, 116, 97, 116, 117, 115, 58, 111, 107, 13, 10, 97, 111, 102, 95, 108, 97, 115, 116, 95, 119, 114, 105, 116, 101, 95, 115, 116, 97, 116, 117, 115, 58, 111, 107, 13, 10, 97, 111, 102, 95, 99, 117, 114, 114, 101, 110, 116, 95, 115, 105, 122, 101, 58, 55, 57, 13, 10, 97, 111, 102, 95, 98, 97, 115, 101, 95, 115, 105, 122, 101, 58, 48, 13, 10, 97, 111, 102, 95, 112, 101, 110, 100, 105, 110, 103, 95, 114, 101, 119, 114, 105, 116, 101, 58, 48, 13, 10, 97, 111, 102, 95, 98, 117, 102, 102, 101, 114, 95, 108, 101, 110, 103, 116, 104, 58, 48, 13, 10, 97, 111, 102, 95, 114, 101, 119, 114, 105, 116, 101, 95, 98, 117, 102, 102, 101, 114, 95, 108, 101, 110, 103, 116, 104, 58, 48, 13, 10, 97, 111, 102, 95, 112, 101, 110, 100, 105, 110, 103, 95, 98, 105, 111, 95, 102, 115, 121, 110, 99, 58, 48, 13, 10, 97, 111, 102, 95, 100, 101, 108, 97, 121, 101, 100, 95, 102, 115, 121, 110, 99, 58, 48, 13, 10, 13, 10, 35, 32, 83, 116, 97, 116, 115, 13, 10, 116, 111, 116, 97, 108, 95, 99, 111, 110, 110, 101, 99, 116, 105, 111, 110, 115, 95, 114, 101, 99, 101, 105, 118, 101, 100, 58, 49, 53, 53, 13, 10, 116, 111, 116, 97, 108, 95, 99, 111, 109, 109, 97, 110, 100, 115, 95, 112, 114, 111, 99, 101, 115, 115, 101, 100, 58, 54, 56, 50, 56, 49, 48, 51, 13, 10, 105, 110, 115, 116, 97, 110, 116, 97, 110, 101, 111, 117, 115, 95, 111, 112, 115, 95, 112, 101, 114, 95, 115, 101, 99, 58, 50, 50, 13, 10, 116, 111, 116, 97, 108, 95, 110, 101, 116, 95, 105, 110, 112, 117, 116, 95, 98, 121, 116, 101, 115, 58, 49, 49, 48, 57, 49, 50, 52, 54, 49, 13, 10, 116, 111, 116, 97, 108, 95, 110, 101, 116, 95, 111, 117, 116, 112, 117, 116, 95, 98, 121, 116, 101, 115, 58, 49, 52, 53, 54, 50, 56, 57, 52, 48, 53, 13, 10, 105, 110, 115, 116, 97, 110, 116, 97, 110, 101, 111, 117, 115, 95, 105, 110, 112, 117, 116, 95, 107, 98, 112, 115, 58, 48, 46, 51, 52, 13, 10, 105, 110, 115, 116, 97, 110, 116, 97, 110, 101, 111, 117, 115, 95, 111, 117, 116, 112, 117, 116, 95, 107, 98, 112, 115, 58, 51, 46, 49, 49, 13, 10, 114, 101, 106, 101, 99, 116, 101, 100, 95, 99, 111, 110, 110, 101, 99, 116, 105, 111, 110, 115, 58, 48, 13, 10, 115, 121, 110, 99, 95, 102, 117, 108, 108, 58, 48, 13, 10, 115, 121, 110, 99, 95, 112, 97, 114, 116, 105, 97, 108, 95, 111, 107, 58, 48, 13, 10, 115, 121, 110, 99, 95, 112, 97, 114, 116, 105, 97, 108, 95, 101, 114, 114, 58, 48, 13, 10, 101, 120, 112, 105, 114, 101, 100, 95, 107, 101, 121, 115, 58, 48, 13, 10, 101, 118, 105, 99, 116, 101, 100, 95, 107, 101, 121, 115, 58, 48, 13, 10, 107, 101, 121, 115, 112, 97, 99, 101, 95, 104, 105, 116, 115, 58, 48, 13, 10, 107, 101, 121, 115, 112, 97, 99, 101, 95, 109, 105, 115, 115, 101, 115, 58, 48, 13, 10, 112, 117, 98, 115, 117, 98, 95, 99, 104, 97, 110, 110, 101, 108, 115, 58, 48, 13, 10, 112, 117, 98, 115, 117, 98, 95, 112, 97, 116, 116, 101, 114, 110, 115, 58, 48, 13, 10, 108, 97, 116, 101, 115, 116, 95, 102, 111, 114, 107, 95, 117, 115, 101, 99, 58, 50, 56, 49, 13, 10, 109, 105, 103, 114, 97, 116, 101, 95, 99, 97, 99, 104, 101, 100, 95, 115, 111, 99, 107, 101, 116, 115, 58, 48, 13, 10, 13, 10, 35, 32, 82, 101, 112, 108, 105, 99, 97, 116, 105, 111, 110, 13, 10, 114, 111, 108, 101, 58, 109, 97, 115, 116, 101, 114, 13, 10, 99, 111, 110, 110, 101, 99, 116, 101, 100, 95, 115, 108, 97, 118, 101, 115, 58, 48, 13, 10, 109, 97, 115, 116, 101, 114, 95, 114, 101, 112, 108, 95, 111, 102, 102, 115, 101, 116, 58, 48, 13, 10, 114, 101, 112, 108, 95, 98, 97, 99, 107, 108, 111, 103, 95, 97, 99, 116, 105, 118, 101, 58, 48, 13, 10, 114, 101, 112, 108, 95, 98, 97, 99, 107, 108, 111, 103, 95, 115, 105, 122, 101, 58, 49, 48, 52, 56, 53, 55, 54, 13, 10, 114, 101, 112, 108, 95, 98, 97, 99, 107, 108, 111, 103, 95, 102, 105, 114, 115, 116, 95, 98, 121, 116, 101, 95, 111, 102, 102, 115, 101, 116, 58, 48, 13, 10, 114, 101, 112, 108, 95, 98, 97, 99, 107, 108, 111, 103, 95, 104, 105, 115, 116, 108, 101, 110, 58, 48, 13, 10, 13, 10, 35, 32, 67, 80, 85, 13, 10, 117, 115, 101, 100, 95, 99, 112, 117, 95, 115, 121, 115, 58, 50, 56, 54, 46, 53, 57, 13, 10, 117, 115, 101, 100, 95, 99, 112, 117, 95, 117, 115, 101, 114, 58, 50, 51, 56, 46, 49, 57, 13, 10, 117, 115, 101, 100, 95, 99, 112, 117, 95, 115, 121, 115, 95, 99, 104, 105, 108, 100, 114, 101, 110, 58, 48, 46, 48, 48, 13, 10, 117, 115, 101, 100, 95, 99, 112, 117, 95, 117, 115, 101, 114, 95, 99, 104, 105, 108, 100, 114, 101, 110, 58, 48, 46, 48, 48, 13, 10, 13, 10, 35, 32, 67, 108, 117, 115, 116, 101, 114, 13, 10, 99, 108, 117, 115, 116, 101, 114, 95, 101, 110, 97, 98, 108, 101, 100, 58, 48, 13, 10, 13, 10, 35, 32, 75, 101, 121, 115, 112, 97, 99, 101, 13, 10, 100, 98, 48, 58, 107, 101, 121, 115, 61, 49, 44, 101, 120, 112, 105, 114, 101, 115, 61, 48, 44, 97, 118, 103, 95, 116, 116, 108, 61, 48, 13, 10}
	var nr int
	{
		fmt.Printf("TestCutRedisInfoSegment case %d.\n", nr)
		nr++

		expect := []byte{35, 32, 83, 101, 114, 118, 101, 114, 13, 10, 114, 101, 100, 105, 115, 95, 118, 101, 114, 115, 105, 111, 110, 58, 51, 46, 50, 46, 49, 49, 13, 10, 114, 101, 100, 105, 115, 95, 103, 105, 116, 95, 115, 104, 97, 49, 58, 100, 101, 49, 97, 100, 48, 50, 54, 13, 10, 114, 101, 100, 105, 115, 95, 103, 105, 116, 95, 100, 105, 114, 116, 121, 58, 48, 13, 10, 114, 101, 100, 105, 115, 95, 98, 117, 105, 108, 100, 95, 105, 100, 58, 54, 101, 49, 99, 101, 51, 55, 97, 102, 101, 48, 54, 99, 51, 52, 99, 13, 10, 114, 101, 100, 105, 115, 95, 109, 111, 100, 101, 58, 115, 116, 97, 110, 100, 97, 108, 111, 110, 101, 13, 10, 111, 115, 58, 76, 105, 110, 117, 120, 32, 51, 46, 49, 48, 46, 48, 45, 57, 53, 55, 46, 50, 49, 46, 51, 46, 101, 108, 55, 46, 120, 56, 54, 95, 54, 52, 32, 120, 56, 54, 95, 54, 52, 13, 10, 97, 114, 99, 104, 95, 98, 105, 116, 115, 58, 54, 52, 13, 10, 109, 117, 108, 116, 105, 112, 108, 101, 120, 105, 110, 103, 95, 97, 112, 105, 58, 101, 112, 111, 108, 108, 13, 10, 103, 99, 99, 95, 118, 101, 114, 115, 105, 111, 110, 58, 52, 46, 56, 46, 53, 13, 10, 112, 114, 111, 99, 101, 115, 115, 95, 105, 100, 58, 57, 48, 54, 57, 13, 10, 114, 117, 110, 95, 105, 100, 58, 49, 52, 49, 98, 51, 54, 101, 102, 50, 51, 50, 57, 56, 52, 53, 54, 98, 52, 100, 97, 101, 48, 49, 51, 49, 50, 57, 49, 98, 50, 100, 101, 99, 99, 101, 102, 102, 52, 54, 53, 13, 10, 116, 99, 112, 95, 112, 111, 114, 116, 58, 56, 53, 48, 51, 13, 10, 117, 112, 116, 105, 109, 101, 95, 105, 110, 95, 115, 101, 99, 111, 110, 100, 115, 58, 53, 56, 57, 54, 51, 49, 13, 10, 117, 112, 116, 105, 109, 101, 95, 105, 110, 95, 100, 97, 121, 115, 58, 54, 13, 10, 104, 122, 58, 49, 48, 13, 10, 108, 114, 117, 95, 99, 108, 111, 99, 107, 58, 51, 55, 52, 50, 48, 50, 48, 13, 10, 101, 120, 101, 99, 117, 116, 97, 98, 108, 101, 58, 47, 100, 97, 116, 97, 47, 109, 97, 112, 103, 111, 111, 47, 99, 111, 100, 105, 115, 47, 46, 47, 98, 105, 110, 47, 99, 111, 100, 105, 115, 45, 115, 101, 114, 118, 101, 114, 13, 10, 99, 111, 110, 102, 105, 103, 95, 102, 105, 108, 101, 58, 47, 100, 97, 116, 97, 47, 109, 97, 112, 103, 111, 111, 47, 99, 111, 100, 105, 115, 47, 46, 47, 99, 111, 100, 105, 115, 95, 115, 101, 114, 118, 101, 114, 47, 99, 111, 110, 102, 47, 114, 101, 100, 105, 115, 45, 56, 53, 48, 51, 46, 99, 111, 110, 102, 13, 10}
		ret := CutRedisInfoSegment(input, "server")
		assert.Equal(t, len(expect), len(ret), "should be equal")
		assert.Equal(t, expect, ret, "should be equal")

		ret = CutRedisInfoSegment(input, "SeRver")
		assert.Equal(t, len(expect), len(ret), "should be equal")
		assert.Equal(t, expect, ret, "should be equal")
	}

	{
		fmt.Printf("TestCutRedisInfoSegment case %d.\n", nr)
		nr++

		expect := []byte{35, 32, 67, 108, 105, 101, 110, 116, 115, 13, 10, 99, 111, 110, 110, 101, 99, 116, 101, 100, 95, 99, 108, 105, 101, 110, 116, 115, 58, 52, 57, 13, 10, 99, 108, 105, 101, 110, 116, 95, 108, 111, 110, 103, 101, 115, 116, 95, 111, 117, 116, 112, 117, 116, 95, 108, 105, 115, 116, 58, 48, 13, 10, 99, 108, 105, 101, 110, 116, 95, 98, 105, 103, 103, 101, 115, 116, 95, 105, 110, 112, 117, 116, 95, 98, 117, 102, 58, 48, 13, 10, 98, 108, 111, 99, 107, 101, 100, 95, 99, 108, 105, 101, 110, 116, 115, 58, 48, 13, 10}
		ret := CutRedisInfoSegment(input, "Clients")
		assert.Equal(t, len(expect), len(ret), "should be equal")
		assert.Equal(t, expect, ret, "should be equal")
	}
}

func TestCompareVersion(t *testing.T) {
	var nr int
	{
		fmt.Printf("TestCompareVersion case %d.\n", nr)
		nr++

		assert.Equal(t, 1, CompareVersion("1.2", "1.3", 2), "should be equal")
		assert.Equal(t, 0, CompareVersion("1.2", "1.3", 1), "should be equal")
		assert.Equal(t, 2, CompareVersion("1.4", "1.3", 2), "should be equal")
		assert.Equal(t, 2, CompareVersion("1.4.x", "1.3", 2), "should be equal")
		assert.Equal(t, 3, CompareVersion("1.4.x", "1.4", 3), "should be equal")
		assert.Equal(t, 0, CompareVersion("1.4.x", "1.4", 0), "should be equal")
		assert.Equal(t, 2, CompareVersion("2.4", "1.1", 2), "should be equal")
	}
}


func TestCompareUnorderedList(t *testing.T) {
	var nr int
	{
		fmt.Printf("TestCompareUnorderedList case %d.\n", nr)
		nr++

		var a, b []string
		assert.Equal(t, true, CompareUnorderedList(a, b), "should be equal")

		a = []string{"1", "2", "3"}
		b = []string{"1", "2"}
		assert.Equal(t, false, CompareUnorderedList(a, b), "should be equal")

		a = []string{"1", "2", "3"}
		b = []string{"3", "1", "2"}
		assert.Equal(t, true, CompareUnorderedList(a, b), "should be equal")

		a = []string{"1", "2", "3"}
		b = []string{"4", "1", "2"}
		assert.Equal(t, false, CompareUnorderedList(a, b), "should be equal")
	}
}

func TestGetSlotDistribution(t *testing.T) {
	var nr int
	{
		fmt.Printf("TestGetSlotDistribution case %d.\n", nr)
		nr++

		ret, err := GetSlotDistribution(testAddrCluster, "auth", "", false)
		assert.Equal(t, nil, err, "should be equal")
		assert.NotEqual(t, 0, len(ret), "should be equal")
		assert.Equal(t, 16383, ret[len(ret) - 1].SlotRightBoundary, "should be equal")
		for i := 1; i < len(ret) - 1; i++ {
			assert.Equal(t, ret[i - 1].SlotRightBoundary + 1, ret[i].SlotLeftBoundary, "should be equal")
			assert.NotEqual(t, "", ret[i].Master, "should be equal")
		}
		fmt.Println(ret)
	}
}

func TestChoseSlotInRange(t *testing.T) {
	var nr int
	{
		fmt.Printf("TestChoseSlotInRange case %d.\n", nr)
		nr++

		ret := ChoseSlotInRange(CheckpointKey, 0, 0)
		fmt.Println(ret)
		assert.NotEqual(t, "", ret, "should be equal")
	}

	{
		fmt.Printf("TestChoseSlotInRange case %d.\n", nr)
		nr++

		// test all slots
		mp := make(map[string]int)
		for i := 0; i < 16384; i++ {
			ret := ChoseSlotInRange(CheckpointKey, i, i)
			assert.NotEqual(t, "", ret, "should be equal")

			fmt.Printf("slot[%v] -> %v\n", i, ret)
			_, ok := mp[ret]
			assert.Equal(t, false, ok, "should be equal")
			mp[ret] = i
		}
		// fmt.Println(mp)
	}
}