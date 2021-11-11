package routers

import (
	"envelope_rain_group10/logger"
	"envelope_rain_group10/nsqclient"
	"envelope_rain_group10/redisclient"
	"envelope_rain_group10/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"time"
)

func LoadSnatch(e *gin.Engine) {
	e.POST("/snatch", SnatchHandler)
}

func SnatchHandler(c *gin.Context) {
	//每个人能抢的最大红包数应该从配置文件读进来

	uidString, _ := c.GetPostForm("uid")
	logger.Logger.Info("snatching envelope", zap.String("uid", uidString))
	uid, _ := strconv.ParseInt(uidString, 10, 64)

	exist, err := redisClient.RedisClient.ExistUser(uid)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	if exist == false {
		err := redisClient.RedisClient.CreateUserInRedis(uid)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
	}

	flag := true
	count, err := redisClient.RedisClient.GetCountWithNextRedPacketByUserId(uid)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	if count == -1 {
		flag = false
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "fail! snatch too much",
			"data": gin.H{
				"max_count": utils.MaxTimes,
				"cur_count": count,
			},
		})
		return
	}

	//根据概率计算用户这次应不应该拿到红包，这里我想的是对所有的请求做统一的处理，直接放弃一部分请求不处理，
	//这样既满足了概率也减轻了后端的压力,只处理十分之一的请求
	rand.Seed(time.Now().Unix())
	probability := int(utils.Probability * 1000000)
	randNum := rand.Intn(1000000)

	if randNum >= probability {
		flag = false
		err := redisClient.RedisClient.ReduceUserGetRedPacketCount(uid)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
		c.JSON(200, gin.H{
			"code": -3,
			"msg":  "According to the probability, the red envelope can not be snatched this time",
		})
		return
	}

	if flag {
		timeStamp := time.Now().Unix()
		redPacketId, err := redisClient.RedisClient.GetRedPacket()
		if err != nil {
			logger.Logger.Error(err.Error())
		}

		if redPacketId == -1 {
			err := redisClient.RedisClient.ReduceUserGetRedPacketCount(uid)
			if err != nil {
				logger.Logger.Error(err.Error())
			}
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "fail! snatch too much",
				"data": gin.H{
					"max_count": utils.MaxTimes,
					"cur_count": utils.MaxTimes,
				},
			})
			return
		}

		nsqclient.ProduceMessage("UpdateCountByUid", strconv.FormatInt(uid, 10))
		value, _ := redisClient.RedisClient.GetRedPacketMoney(redPacketId)
		nsqclient.ProduceMessage("CreateEnvelopeDetail", fmt.Sprintf("%d,%d,%d,%d", redPacketId, uid, value, timeStamp))

		err = redisClient.RedisClient.AddToUserRedPacketList(uid, redPacketId)
		if err != nil {
			logger.Logger.Error(err.Error())
		}

		err = redisClient.RedisClient.AddToUserRedPacketTimeList(uid, timeStamp)
		if err != nil {
			logger.Logger.Error(err.Error())
		}

		err = redisClient.RedisClient.MakeWalletCacheInvalid(uid)
		if err != nil {
			logger.Logger.Error(err.Error())
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"envelope_id": redPacketId,
				"max_count":   utils.MaxTimes,
				"cur_count":   count,
			},
		})
		return
	}
}
