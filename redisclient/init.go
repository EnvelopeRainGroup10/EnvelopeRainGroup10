package redisClient

import (
	"envelope_rain_group10/logger"
	"envelope_rain_group10/utils"
)

var RedisClient *redisClient

func InitRedisClient() {
	// var addr = "redis-cn02a2vagk7trou1z-direct.redis.ivolces.com:6380"
	var addr = "redis-cn02a2vagk7trou1z.redis.volces.com:6379"
	var password = "Njuse2021"
	var db int64 = 0
	var poolSize int64 = 10000
	var maxPacketNum = utils.TotalNum
	var maxGetNum = utils.MaxTimes
	var keyPre = "test1:"
	r, err := NewRedisClient(addr, password, db, poolSize, maxPacketNum, maxGetNum, keyPre)
	RedisClient = r
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
