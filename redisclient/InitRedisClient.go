package redisClient

import (
	"envelope_rain_group10/logger"
	"envelope_rain_group10/utils"
)

var RedisClient *redisClient

func InitRedisClient()  {
	var addr string = "127.0.0.1:6379"
	var password string = ""
	var db int64 = 0
	var poolSize int64 = 1000
	var maxPacketNum int64 = utils.TotalNum
	var maxGetNum int64 = utils.MaxTimes
	var keyPre string = "test1:"
	r, err := NewRedisClient(addr, password, db, poolSize, maxPacketNum, maxGetNum, keyPre)
	RedisClient = r
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
