package routers

import (
	"envelope_rain_group10/redisclient"
	"envelope_rain_group10/sql"
	"envelope_rain_group10/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"strconv"
	"time"
)


func LoadSnatch(e *gin.Engine) {
	e.POST("/snatch", SnatchHandler)
}

func SnatchHandler(c *gin.Context) {
	//每个人能抢的最大红包数应该从配置文件读进来

	uid, _ := c.GetPostForm("uid")
	logs.Printf("%s is snatching envelope", uid)
	int_uid, _ := strconv.ParseInt(uid, 10, 64)

	exist, err := redisClient.RedisClient.ExistUser(int_uid)
	if err != nil {
		log.Println(err)
	}

	if exist == false {
		sql.GetUser(int_uid)
		err := redisClient.RedisClient.CreateUserInRedis(int_uid)
		if err != nil {
			logs.Println(err)
		}
	}
	flag := true
	count, err := redisClient.RedisClient.GetCountWithNextRedPacketByUserId(int_uid)
	if err != nil {
		logs.Println(err)
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
	rand_num := rand.Intn(1000000)

	fmt.Println(rand_num)
	fmt.Println(probability)

	if rand_num >= probability {
		flag = false
		redisClient.RedisClient.ReduceUserGetRedPacketCount(int_uid)
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
			logs.Println(err)
		}

		if redPacketId==-1{
			redisClient.RedisClient.ReduceUserGetRedPacketCount(int_uid)//没抢到，还原redis中用户抢到的红包计数
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

		sql.UpdateCountByUid(int_uid)
		value, _ := redisClient.RedisClient.GetRedPacketMoney(redPacketId)
		sql.CreateEnvelopeDetail(redPacketId, int_uid, value, timeStamp)

		redisClient.RedisClient.AddToUserRedPacketList(int_uid, redPacketId)
		redisClient.RedisClient.AddToUserRedPacketTimeList(int_uid,timeStamp)
		redisClient.RedisClient.MakeWalletCacheInvalid(int_uid)
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
