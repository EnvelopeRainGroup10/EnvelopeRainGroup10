package redisClient

import (
	"envelope_rain_group10/logger"
	"envelope_rain_group10/utils"
)

var RedisClient *redisClient

func InitRedisClient()  {
	var addr = "127.0.0.1:6379"
	var password = ""
	var db int64 = 0
	var poolSize int64 = 1000
	var maxPacketNum = utils.TotalNum
	var maxGetNum = utils.MaxTimes
	var keyPre = "test1:"
	r, err := NewRedisClient(addr, password, db, poolSize, maxPacketNum, maxGetNum, keyPre)
	RedisClient = r
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
